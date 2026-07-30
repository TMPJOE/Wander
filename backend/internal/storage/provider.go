// Package storage abstracts file uploads behind a single Provider interface
// so the upload handler doesn't care where the bytes end up.
//
// Two implementations ship in this package:
//
//   - LocalProvider: writes to a directory on disk and returns a server-relative
//     URL (e.g. "/uploads/abc.jpg") that the Go server then serves via
//     http.FileServer. Matches the original behavior.
//
//   - S3Provider: streams bytes into any S3-protocol-compatible bucket
//     (AWS S3, Cloudflare R2, MinIO, Supabase, Backblaze B2, ...). Returns
//     absolute bucket URLs so the Go server no longer needs to serve the
//     /uploads/* route.
//
// Selection happens in config.Storage and is wired in cmd/server/main.go.
package storage

import (
	"context"
	"io"
)

// SaveResult is what a provider returns after persisting a file.
//
// URL is the value that gets handed back to the client and stored in the DB.
// For local mode it is server-relative ("/uploads/abc.jpg"). For S3 mode it
// is an absolute URL ("https://cdn.example.com/abc.jpg").
//
// Key is the provider-internal identifier of the stored object (filename for
// local, object key for S3). Returned for logging and potential future
// deletion flows.
type SaveResult struct {
	URL string
	Key string
}

// Provider is the upload abstraction used by handler.UploadHandler.
type Provider interface {
	// Save reads up to limit bytes from r and stores them under ext, tagged
	// with contentType. Returns the public URL and object key on success.
	// The caller is responsible for closing r after Save returns.
	Save(ctx context.Context, r io.Reader, contentType, ext string) (SaveResult, error)

	// Delete removes a previously stored object identified by its public URL.
	// It derives the internal provider key from the URL. A missing object is
	// treated as success (idempotent) so callers can clean up best-effort
	// without surfacing "not found" errors from the backend.
	Delete(ctx context.Context, url string) error

	// DeleteMany removes a batch of objects in one call. Each URL is resolved
	// the same way as Delete. Errors for individual objects are logged but do
	// not abort the rest of the batch; the returned error is nil unless the
	// batch itself could not be initiated.
	DeleteMany(ctx context.Context, urls []string) error
}
