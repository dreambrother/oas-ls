package handler

import (
	"testing"

	"github.com/owenrumney/go-lsp/servertest"
)

func TestHover(t *testing.T) {
    h := servertest.New(t, &Handler{})

    h.DidOpen("file:///test.txt", "plaintext", "hello")
    hover, err := h.Hover("file:///test.txt", 0, 0)
    if err != nil {
        t.Fatal(err)
    }
    if hover == nil {
        t.Fatal("expected hover result")
    }
}
