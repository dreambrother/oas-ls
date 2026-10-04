package handler

import (
	"bytes"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func newTestServer(t *testing.T) *servertest.Harness {
	return servertest.New(t, NewHandler(newTestLogger(t)))
}

func newTestLogger(t *testing.T) *slog.Logger {
	t.Helper()
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("captured logs:\n%s", buf.String())
		}
	})
	return log
}

func loadTestdata(t *testing.T, name string) (lsp.DocumentURI, string) {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return lsp.DocumentURI(toFileURI(path)), string(data)
}

func toFileURI(path string) string {
	return (&url.URL{Scheme: "file", Path: filepath.ToSlash(path)}).String()
}
