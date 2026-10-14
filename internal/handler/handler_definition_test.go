package handler

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/dreambrother/oas-ls/internal/file"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func TestDefinition(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")

	h.DidOpen(uri, "yaml", content)
	links := definitionLinks(t, h, uri, 17, 39) // components.yaml#/CreateResult

	expectDefinition(t, links, "definition/internal/components.yaml", 6, 0, 6, 12)
	expectOriginRange(t, links, 17, 19, 17, 50)
}

func TestDefinition_AfterDocChange(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")
	_, changedContent := loadTestdata(t, "definition/internal/create-changed.yaml")

	h.DidOpen(uri, "yaml", content)
	h.DidChange(uri, 2, changedContent)
	links := definitionLinks(t, h, uri, 10, 34) // components.yaml#/CreateRequest2

	expectDefinition(t, links, "definition/internal/components.yaml", 12, 0, 12, 14)
	expectOriginRange(t, links, 10, 17, 10, 48)
}

func TestDefinition_SameDocComponentsSchemas(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/search.yaml")

	h.DidOpen(uri, "yaml", content)
	links := definitionLinks(t, h, uri, 17, 40) // #/components/schemas/SearchResult

	expectDefinition(t, links, "definition/internal/search.yaml", 23, 4, 23, 16)
	expectOriginRange(t, links, 17, 19, 17, 52)
}

func TestDefinition_SameDocRef(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/components.yaml")

	h.DidOpen(uri, "yaml", content)
	links := definitionLinks(t, h, uri, 10, 15) // #/SomeType

	expectDefinition(t, links, "definition/internal/components.yaml", 15, 0, 15, 8)
	expectOriginRange(t, links, 10, 13, 10, 23)
}

func TestDefinition_FileRef(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/api.yaml")

	h.DidOpen(uri, "yaml", content)
	links := definitionLinks(t, h, uri, 3, 20) // internal/create.yaml

	expectDefinition(t, links, "definition/internal/create.yaml", 0, 0, 0, 0)
	expectOriginRange(t, links, 3, 11, 3, 31)
}

func expectDefinition(
	t *testing.T,
	links []lsp.LocationLink,
	path string,
	startLine, startCharacter, endLine, endCharacter int,
) {
	t.Helper()
	if len(links) != 1 {
		t.Fatal("expected 1 definition but was", len(links))
	}
	d := links[0]
	expectedPath, _ := filepath.Abs("testdata/" + path)
	if expected := file.ToURI(expectedPath); string(d.TargetURI) != expected {
		t.Errorf("expected definition URI to be '%s', but was '%s'", expected, d.TargetURI)
	}
	if d.TargetRange.Start.Line != startLine {
		t.Errorf("expected definition start line to be '%d', but was '%d'", startLine, d.TargetRange.Start.Line)
	}
	if d.TargetRange.Start.Character != startCharacter {
		t.Errorf("expected definition start character to be '%d', but was '%d'", startCharacter, d.TargetRange.Start.Character)
	}
	if d.TargetRange.End.Line != endLine {
		t.Errorf("expected definition end line to be '%d', but was '%d'", endLine, d.TargetRange.End.Line)
	}
	if d.TargetRange.End.Character != endCharacter {
		t.Errorf("expected definition end character to be '%d', but was '%d'", endCharacter, d.TargetRange.End.Character)
	}
}

func expectOriginRange(
	t *testing.T,
	links []lsp.LocationLink,
	startLine, startCharacter, endLine, endCharacter int,
) {
	t.Helper()
	if len(links) != 1 {
		t.Fatalf("expected 1 definition but was %d", len(links))
	}
	origin := links[0].OriginSelectionRange
	if origin == nil {
		t.Fatal("OriginSelectionRange is nil")
	}
	if origin.Start.Line != startLine {
		t.Errorf("expected origin start line to be '%d', but was '%d'", startLine, origin.Start.Line)
	}
	if origin.Start.Character != startCharacter {
		t.Errorf("expected origin start character to be '%d', but was '%d'", startCharacter, origin.Start.Character)
	}
	if origin.End.Line != endLine {
		t.Errorf("expected origin end line to be '%d', but was '%d'", endLine, origin.End.Line)
	}
	if origin.End.Character != endCharacter {
		t.Errorf("expected origin end character to be '%d', but was '%d'", endCharacter, origin.End.Character)
	}
}

func TestDefinition_NotRef(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/api.yaml")

	h.DidOpen(uri, "yaml", content)
	links := definitionLinks(t, h, uri, 1, 1)

	if len(links) != 0 {
		t.Fatal("expected 0 definition but was", len(links))
	}
}

func TestDefinition_UnknownFile(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/components.yaml")

	h.DidOpen(uri, "yaml", content)
	links := definitionLinks(t, h, uri, 4, 28)

	if len(links) != 0 {
		t.Fatal("expected 0 definition but was", len(links))
	}
}

func definitionLinks(t *testing.T, h *servertest.Harness, uri lsp.DocumentURI, line, char int) []lsp.LocationLink {
	t.Helper()
	raw, err := h.Call("textDocument/definition", &lsp.DefinitionParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Position:     lsp.Position{Line: line, Character: char},
	})
	if err != nil {
		t.Fatal(err)
	}
	var links []lsp.LocationLink
	if err := json.Unmarshal(raw, &links); err != nil {
		t.Fatalf("unmarshal definitions: %v", err)
	}
	return links
}
