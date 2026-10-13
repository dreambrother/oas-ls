package file

import (
	"net/url"
	"path/filepath"

	"github.com/owenrumney/go-lsp/lsp"
)

func ToDocumentURI(path string) lsp.DocumentURI {
	return lsp.DocumentURI(ToURI(path))
}

func ToURI(path string) string {
	return (&url.URL{Scheme: "file", Path: filepath.ToSlash(path)}).String()
}
