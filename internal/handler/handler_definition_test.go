package handler

import (
	"path/filepath"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
)

func TestDefinition(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 17, 39) // components.yaml#/CreateResult

	if err != nil {
		t.Fatal(err)
	}
	expectDefinition(t, definitions, "definition/internal/components.yaml", 6, 0, 6, 0)
}

func TestDefinition_AfterDocChange(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")
	_, changedContent := loadTestdata(t, "definition/internal/create-changed.yaml")

	h.DidOpen(uri, "yaml", content)
	h.DidChange(uri, 2, changedContent)
	definitions, err := h.Definition(uri, 10, 34) // components.yaml#/CreateRequest2

	if err != nil {
		t.Fatal(err)
	}
	expectDefinition(t, definitions, "definition/internal/components.yaml", 12, 0, 12, 0)
}

func TestDefinition_SameDocComponentsSchemas(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/search.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 17, 40) // #/components/schemas/SearchResult

	if err != nil {
		t.Fatal(err)
	}
	expectDefinition(t, definitions, "definition/internal/search.yaml", 23, 4, 23, 4)
}

func TestDefinition_SameDocRef(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/components.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 10, 15) // #/SomeType

	if err != nil {
		t.Fatal(err)
	}
	expectDefinition(t, definitions, "definition/internal/components.yaml", 15, 0, 15, 0)
}

func expectDefinition(
	t *testing.T,
	definitions []lsp.Location,
	path string,
	startLine, startCharacter, endLine, endCharacter int,
) {
	if len(definitions) != 1 {
		t.Fatal("expected 1 definition but was", len(definitions))
	}
	d := definitions[0]
	expectedPath, _ := filepath.Abs("testdata/" + path)
	if expected := toFileURI(expectedPath); string(d.URI) != expected {
		t.Errorf("expected definition URI to be '%s', but was '%s'", expected, d.URI)
	}
	if d.Range.Start.Line != startLine {
		t.Errorf("expected definition start line to be '%d', but was '%d'", startLine, d.Range.Start.Line)
	}
	if d.Range.Start.Character != startCharacter {
		t.Errorf("expected definition start character to be '%d', but was '%d'", startCharacter, d.Range.Start.Character)
	}
	if d.Range.End.Line != endLine {
		t.Errorf("expected definition end line to be '%d', but was '%d'", endLine, d.Range.End.Line)
	}
	if d.Range.End.Character != endCharacter {
		t.Errorf("expected definition end character to be '%d', but was '%d'", endCharacter, d.Range.End.Character)
	}
}

func TestDefinition_NotRef(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/api.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 1, 1)

	if err != nil {
		t.Fatal(err)
	}
	if len(definitions) != 0 {
		t.Fatal("expected 0 definition but was", len(definitions))
	}
}

func TestDefinition_UnknownFile(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/components.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 4, 28)

	if err != nil {
		t.Fatal(err)
	}
	if len(definitions) != 0 {
		t.Fatal("expected 0 definition but was", len(definitions))
	}
}
