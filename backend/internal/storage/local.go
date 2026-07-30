package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// LocalProvider stores files on local disk and returns a server-relative
// URL the Go server serves via http.FileServer. This is the original
// Wander upload behavior preserved unchanged.
type LocalProvider struct {
	uploadsDir    string // absolute path on disk
	publicBaseURL string // URL prefix returned to clients, e.g. "/uploads"
}

// NewLocalProvider returns a LocalProvider rooted at dir, with files served
// under publicBaseURL. publicBaseURL must be a server-relative path rooted
// at "/" (no trailing slash).
func NewLocalProvider(dir, publicBaseURL string) (*LocalProvider, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, err
	}
	publicBaseURL = strings.TrimRight(publicBaseURL, "/")
	if publicBaseURL == "" {
		publicBaseURL = "/uploads"
	}
	return &LocalProvider{
		uploadsDir:    abs,
		publicBaseURL: publicBaseURL,
	}, nil
}

// Dir returns the absolute on-disk directory. Used by main.go to wire the
// http.FileServer mount point.
func (p *LocalProvider) Dir() string { return p.uploadsDir }

// BaseURL returns the URL prefix files are served under. Used by main.go
// to register the static route at the matching path.
func (p *LocalProvider) BaseURL() string { return p.publicBaseURL }

// Save implements Provider. The contentType argument is ignored for local
// storage — http.FileServer infers the Content-Type from the file extension
// at serve time.
func (p *LocalProvider) Save(_ context.Context, r io.Reader, _, ext string) (SaveResult, error) {
	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return SaveResult{}, err
	}
	key := hex.EncodeToString(randBytes) + ext

	if err := os.MkdirAll(p.uploadsDir, 0o755); err != nil {
		return SaveResult{}, err
	}
	dst, err := os.Create(filepath.Join(p.uploadsDir, key))
	if err != nil {
		return SaveResult{}, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r); err != nil {
		return SaveResult{}, err
	}
	return SaveResult{
		URL: p.publicBaseURL + "/" + key,
		Key: key,
	}, nil
}

// keyFromURL extracts the provider key (filename) from a public URL produced
// by Save. It tolerates both the server-relative form ("/uploads/abc.jpg")
// and any absolute URL that happens to resolve to this provider's base. URLs
// not rooted at the provider's base are assumed to be external and returned
// as-is after stripping a leading slash, so Delete is a safe no-op for them.
func (p *LocalProvider) keyFromURL(url string) string {
	if url == "" {
		return ""
	}
	// Strip scheme/host if an absolute URL was passed in by mistake.
	if idx := strings.Index(url, "://"); idx >= 0 {
		if rest := strings.IndexByte(url[idx+3:], '/'); rest >= 0 {
			url = url[idx+3+rest:]
		} else {
			return ""
		}
	}
	if strings.HasPrefix(url, p.publicBaseURL+"/") {
		return strings.TrimPrefix(url, p.publicBaseURL+"/")
	}
	return strings.TrimPrefix(url, "/")
}

// Delete implements Provider. It removes the on-disk file for the given URL.
// A missing file is treated as success so callers can clean up best-effort.
func (p *LocalProvider) Delete(_ context.Context, url string) error {
	key := p.keyFromURL(url)
	if key == "" || strings.ContainsAny(key, `\/`) {
		// Empty key or something that escaped the uploads dir — skip rather
		// than risk touching an arbitrary path.
		return nil
	}
	path := filepath.Join(p.uploadsDir, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// DeleteMany implements Provider. Each URL is deleted independently; a failure
// on one object is logged but does not abort the rest of the batch.
func (p *LocalProvider) DeleteMany(ctx context.Context, urls []string) error {
	for _, url := range urls {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := p.Delete(ctx, url); err != nil {
			slog.Warn("local delete failed", "url", url, "error", err)
		}
	}
	return nil
}
