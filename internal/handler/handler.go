package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dreambrother/oas-ls/internal/operations"
	"github.com/dreambrother/oas-ls/internal/textpos"
	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
)

type Handler struct {
	documents *document.Store
	log       *slog.Logger
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		documents: document.NewStore(),
		log:       log,
	}
}

func (h *Handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	encoding := textpos.PositionEncoding
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{PositionEncoding: &encoding},
		ServerInfo:   &lsp.ServerInfo{Name: "oas-ls", Version: "0.1.0"},
	}, nil
}

func (h *Handler) DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) error {
	_, err := h.documents.Open(params)
	return err
}

func (h *Handler) DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams) error {
	_, err := h.documents.Change(params)
	return err
}

func (h *Handler) DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) error {
	h.documents.Close(params)
	return nil
}

func (h *Handler) Shutdown(_ context.Context) error { return nil }

func (h *Handler) Definition(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.Location, error) {
	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return nil, fmt.Errorf("document not found: %s", params.TextDocument.URI)
	}
	log := h.log.With("operation", "definition")
	return operations.Definition(ctx, log, doc, params)
}
