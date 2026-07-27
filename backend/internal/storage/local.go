package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
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
