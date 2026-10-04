package handler

import (
	"path/filepath"
	"testing"
)

func TestDefinition(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 10, 37)

	if err != nil {
		t.Fatal(err)
	}
	if len(definitions) != 1 {
		t.Fatal("expected 1 definition but was", len(definitions))
	}
	d := definitions[0]
	expectedPath, _ := filepath.Abs("testdata/definition/internal/components.yaml")
	if expected := toFileURI(expectedPath); string(d.URI) != expected {
		t.Errorf("expected definition URI to be '%s', but was '%s'", expected, d.URI)
	}
	if d.Range.Start.Line != 1 {
		t.Error("expected definition start line to be 1, but was", d.Range.Start.Line)
	}
	if d.Range.Start.Character != 1 {
		t.Error("expected definition start character to be 1, but was", d.Range.Start.Character)
	}
	if d.Range.End.Line != 1 {
		t.Error("expected definition end line to be 1, but was", d.Range.End.Line)
	}
	if d.Range.End.Character != 1 {
		t.Error("expected definition end character to be 1, but was", d.Range.End.Character)
	}
}
