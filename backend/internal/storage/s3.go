package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
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
	Bucket        string // required
	Region        string // required, "auto" works for R2/MinIO
	Endpoint      string // custom endpoint for R2/MinIO/Supabase; empty for AWS S3
	AccessKey     string // required
	SecretKey     string // required
	ForcePathStyle bool  // true for MinIO; false for AWS/R2 virtual-hosted
	PublicBaseURL string // optional CDN/bucket URL prefix returned to clients
	KeyPrefix     string // optional prefix for stored keys, e.g. "uploads/"
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
	escaped := url.PathEscape(key)
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", p.bucket, escaped)
}
