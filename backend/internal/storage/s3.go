package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	neturl "net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Provider writes files into any S3-protocol-compatible bucket and returns
// absolute object URLs so the Go server can stop serving /uploads/* routes.
type S3Provider struct {
	client        *s3.Client
	bucket        string
	publicBaseURL string
	keyPrefix     string
}

// S3Options configures NewS3Provider.
type S3Options struct {
	Bucket         string // required
	Region         string // required, "auto" works for R2/MinIO
	Endpoint       string // custom endpoint for R2/MinIO/Supabase; empty for AWS S3
	AccessKey      string // required
	SecretKey      string // required
	ForcePathStyle bool   // true for MinIO; false for AWS/R2 virtual-hosted
	PublicBaseURL  string // optional CDN/bucket URL prefix returned to clients
	KeyPrefix      string // optional prefix for stored keys, e.g. "uploads/"
}

// NewS3Provider builds an S3-compatible client and validates required opts.
// Uses static credentials; region defaults to "auto" for R2/MinIO.
func NewS3Provider(ctx context.Context, opts S3Options) (*S3Provider, error) {
	if opts.Bucket == "" {
		return nil, errors.New("storage: S3_BUCKET is required")
	}
	if opts.AccessKey == "" || opts.SecretKey == "" {
		return nil, errors.New("storage: S3_ACCESS_KEY and S3_SECRET_KEY are required")
	}
	region := opts.Region
	if region == "" {
		region = "auto"
	}

	loadOpts := []func(*awscfg.LoadOptions) error{
		awscfg.WithRegion(region),
		awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			opts.AccessKey, opts.SecretKey, "",
		)),
	}
	cfg, err := awscfg.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if opts.Endpoint != "" {
			o.BaseEndpoint = aws.String(opts.Endpoint)
		}
		o.UsePathStyle = opts.ForcePathStyle
	})

	// If the user gave us a separate public origin (e.g. a Cloudflare CDN in
	// front of the R2 bucket), normalize it. Otherwise fall back to the
	// endpoint URL itself if set, or to the AWS S3 virtual-hosted URL.
	var publicBase string
	if opts.PublicBaseURL != "" {
		publicBase = strings.TrimRight(opts.PublicBaseURL, "/")
	} else if opts.Endpoint != "" {
		// For path-style endpoints (MinIO) the bucket lives under the host;
		// for virtual-hosted endpoints the bucket is in the path here too,
		// since the SDK already does URL rewriting for us when building
		// object URLs — we always fabricate the public URL as host/bucket/key
		// so it works across R2/MinIO/AWS consistently.
		publicBase = strings.TrimRight(opts.Endpoint, "/") + "/" + opts.Bucket
	}
	// AWS S3 with no public base and no endpoint: constructs the standard
	// https://<bucket>.s3.<region>.amazonaws.com/<key> URL below.

	return &S3Provider{
		client:        s3Client,
		bucket:        opts.Bucket,
		publicBaseURL: publicBase,
		keyPrefix:     strings.TrimPrefix(strings.TrimRight(opts.KeyPrefix, "/"), "/"),
	}, nil
}

// Save streams r into the bucket under a random key + ext, tagged with
// contentType. The returned URL is absolute and handed back to the client.
func (p *S3Provider) Save(ctx context.Context, r io.Reader, contentType, ext string) (SaveResult, error) {
	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return SaveResult{}, err
	}
	key := hex.EncodeToString(randBytes) + ext
	if p.keyPrefix != "" {
		key = p.keyPrefix + key
	}

	put := &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
	}
	if _, err := p.client.PutObject(ctx, put); err != nil {
		return SaveResult{}, fmt.Errorf("storage: s3 put: %w", err)
	}

	return SaveResult{URL: p.objectURL(key), Key: key}, nil
}

// keyFromURL extracts the S3 object key from a public URL produced by Save.
// It prefers stripping the configured publicBaseURL prefix (covering CDN and
// custom-endpoint cases). For plain AWS S3 virtual-hosted URLs, it strips the
// known host forms. URLs that do not resolve to this bucket are returned as
// empty so Delete becomes a safe no-op.
func (p *S3Provider) keyFromURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	// Fast path: the URL starts with our configured public base.
	if p.publicBaseURL != "" && strings.HasPrefix(rawURL, p.publicBaseURL+"/") {
		return strings.TrimPrefix(rawURL, p.publicBaseURL+"/")
	}
	// AWS-hosted default form used when no public base was set:
	//   https://<bucket>.s3.amazonaws.com/<key>
	//   https://<bucket>.s3.<region>.amazonaws.com/<key>
	//   https://s3.amazonaws.com/<bucket>/<key>
	parsed, err := neturl.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return ""
	}
	key := strings.TrimPrefix(parsed.Path, "/")
	// Strip a leading "<bucket>/" for the path-form host.
	if keyPrefix := p.bucket + "/"; strings.HasPrefix(key, keyPrefix) {
		key = strings.TrimPrefix(key, keyPrefix)
	}
	// The object key may be percent-escaped in the URL path; decode it so the
	// DeleteObject Key matches what was stored. Unescape errors fall back to
	// the raw segment, which is still correct for the common ASCII filenames.
	if unescaped, err := neturl.PathUnescape(key); err == nil {
		key = unescaped
	}
	return key
}

// Delete implements Provider. It removes the object for the given URL from the
// bucket. A missing object is treated as success (S3 DeleteObject is already
// idempotent, so we just propagate other errors).
func (p *S3Provider) Delete(ctx context.Context, url string) error {
	key := p.keyFromURL(url)
	if key == "" {
		return nil
	}
	if _, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("storage: s3 delete: %w", err)
	}
	return nil
}

// DeleteMany implements Provider. It deletes each object independently; a
// failure on one object is logged but does not abort the rest of the batch.
func (p *S3Provider) DeleteMany(ctx context.Context, urls []string) error {
	for _, url := range urls {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := p.Delete(ctx, url); err != nil {
			slog.Warn("s3 delete failed", "url", url, "error", err)
		}
	}
	return nil
}

// objectURL builds the public URL for a stored object key.
//
// - If publicBaseURL is set (CDN or custom origin), it takes precedence.
// - Otherwise, for path-style/custom endpoints, returns <endpoint>/<bucket>/<key>.
// - Otherwise builds the standard AWS S3 virtual-hosted URL.
func (p *S3Provider) objectURL(key string) string {
	if p.publicBaseURL != "" {
		return p.publicBaseURL + "/" + key
	}
	// Fall through to AWS-hosted style. We assume https:// unless the bucket
	// name gave us a region; this is the common case and matches the SDK's
	// default behavior. Custom endpoints were already handled above where we
	// set publicBaseURL from opts.Endpoint.
	escaped := neturl.PathEscape(key)
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", p.bucket, escaped)
}
