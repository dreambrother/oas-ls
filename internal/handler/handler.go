package handler

import (
	"context"

	"github.com/owenrumney/go-lsp/lsp"
)

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct{}

func (h *Handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "specpls", Version: "0.1.0"},
	}, nil
}

func (h *Handler) Shutdown(_ context.Context) error { return nil }

func (h *Handler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	return &lsp.Hover{
		Contents: lsp.NewHoverContents(lsp.Markdown, "**todo**"),
	}, nil
}
