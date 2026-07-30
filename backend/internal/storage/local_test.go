package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalProvider_SaveDeleteRoundtrip(t *testing.T) {
	tmp := t.TempDir()
	p, err := NewLocalProvider(tmp, "/uploads")
	if err != nil {
		t.Fatalf("NewLocalProvider: %v", err)
	}

	content := []byte("hello world")
	res, err := p.Save(context.Background(), bytes.NewReader(content), "image/jpeg", ".jpg")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	if !strings.HasPrefix(res.URL, "/uploads/") {
		t.Fatalf("URL = %q, want /uploads/ prefix", res.URL)
	}

	path := filepath.Join(tmp, res.Key)
	if got, err := os.ReadFile(path); err != nil {
		t.Fatalf("file not persisted on disk: %v", err)
	} else if !bytes.Equal(got, content) {
		t.Fatalf("disk content mismatch: %q", got)
	}

	if err := p.Delete(context.Background(), res.URL); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("file still on disk after Delete; stat err = %v", err)
	}

	// Delete must be idempotent on a missing object.
	if err := p.Delete(context.Background(), res.URL); err != nil {
		t.Fatalf("idempotent Delete on missing file: %v", err)
	}
}

func TestLocalProvider_DeleteIgnoresExternalURLs(t *testing.T) {
	tmp := t.TempDir()
	p, err := NewLocalProvider(tmp, "/uploads")
	if err != nil {
		t.Fatalf("NewLocalProvider: %v", err)
	}
	// An absolute external URL should not touch any local path and not error.
	if err := p.Delete(context.Background(), "https://example.com/some.jpg"); err != nil {
		t.Fatalf("Delete external URL: %v", err)
	}
	// Empty URL is a no-op.
	if err := p.Delete(context.Background(), ""); err != nil {
		t.Fatalf("Delete empty URL: %v", err)
	}
}
