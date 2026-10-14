package operation

import (
	"context"
	"log/slog"
	"testing"

	"github.com/dreambrother/oas-ls/internal/file"
	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
)

func TestParseRef(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantRef   string
		wantStart int
		wantEnd   int
		wantOK    bool
	}{
		{
			name:      "plain scalar",
			line:      "$ref: ./a.yaml#/Foo",
			wantRef:   "./a.yaml#/Foo",
			wantStart: 6,
			wantEnd:   19,
			wantOK:    true,
		},
		{
			name:      "plain scalar with trailing comment",
			line:      "$ref: ./a.yaml#/Foo # comment",
			wantRef:   "./a.yaml#/Foo",
			wantStart: 6,
			wantEnd:   19,
			wantOK:    true,
		},
		{
			name:      "single-quoted scalar",
			line:      "$ref: './a.yaml#/Foo'",
			wantRef:   "./a.yaml#/Foo",
			wantStart: 7,
			wantEnd:   20,
			wantOK:    true,
		},
		{
			name:      "double-quoted scalar",
			line:      "$ref: \"#/SomeType\"",
			wantRef:   "#/SomeType",
			wantStart: 7,
			wantEnd:   17,
			wantOK:    true,
		},
		{
			name:      "indented single-quoted scalar",
			line:      "  $ref: '#/SomeType'",
			wantRef:   "#/SomeType",
			wantStart: 9,
			wantEnd:   19,
			wantOK:    true,
		},
		{
			name:      "plain scalar with non-ascii",
			line:      "$ref: ./каталог.yaml#/Foo",
			wantRef:   "./каталог.yaml#/Foo",
			wantStart: 6,
			wantEnd:   25,
			wantOK:    true,
		},
		{
			name:   "not a ref",
			line:   "paths:",
			wantOK: false,
		},
		{
			name:   "ref is not the first key",
			line:   "description: $ref",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, start, end, ok := parseRef(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if ref != tt.wantRef {
				t.Errorf("ref = %q, want %q", ref, tt.wantRef)
			}
			if start != tt.wantStart {
				t.Errorf("start = %d, want %d", start, tt.wantStart)
			}
			if end != tt.wantEnd {
				t.Errorf("end = %d, want %d", end, tt.wantEnd)
			}
		})
	}
}

const originTestContent = `
paths:
  /api/internal:
    $ref: '#/components/schemas/Foo'
components:
  schemas:
    Foo:
      type: object
`

func TestDefinition_OriginSelectionRange(t *testing.T) {
	uri, store, doc := openDocument(t, originTestContent)

	links, err := Definition(
		context.Background(),
		slog.New(slog.DiscardHandler),
		doc,
		&lsp.DefinitionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: 3},
		},
		store,
	)
	if err != nil {
		t.Fatalf("Definition: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	link := links[0]

	if link.TargetURI != uri {
		t.Errorf("TargetURI = %q, want %q", link.TargetURI, uri)
	}

	if link.OriginSelectionRange == nil {
		t.Fatal("OriginSelectionRange is nil")
	}
	if got, want := *link.OriginSelectionRange, (lsp.Range{
		Start: lsp.Position{Line: 3, Character: 11}, // after the opening quote
		End:   lsp.Position{Line: 3, Character: 35}, // '#/components/schemas/Foo' is 24 runes long
	}); got != want {
		t.Errorf("OriginSelectionRange = %+v, want %+v", got, want)
	}

	if got, want := link.TargetRange, (lsp.Range{
		Start: lsp.Position{Line: 6, Character: 4},
		End:   lsp.Position{Line: 6, Character: 7}, // 'Foo' is 3 runes long
	}); got != want {
		t.Errorf("TargetRange = %+v, want %+v", got, want)
	}
}

func TestDefinition_NotRef(t *testing.T) {
	uri, store, doc := openDocument(t, originTestContent)

	links, err := Definition(
		context.Background(),
		slog.New(slog.DiscardHandler),
		doc,
		&lsp.DefinitionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: 1},
		},
		store,
	)
	if err != nil {
		t.Fatalf("Definition: %v", err)
	}
	if len(links) != 0 {
		t.Fatalf("expected 0 links, got %d", len(links))
	}
}

func openDocument(t *testing.T, content string) (lsp.DocumentURI, *document.Store, *document.Document) {
	t.Helper()

	uri := file.ToDocumentURI("/virtual/search.yaml")
	store := document.NewStore()
	if _, err := store.Open(&lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "yaml",
			Version:    1,
			Text:       content,
		},
	}); err != nil {
		t.Fatalf("open document: %v", err)
	}
	doc, ok := store.Get(uri)
	if !ok {
		t.Fatal("document not found")
	}
	return uri, store, doc
}
