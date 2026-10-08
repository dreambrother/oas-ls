package handler

import (
	"path/filepath"
	"testing"
)

func TestDefinition(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")

	h.DidOpen(uri, "yaml", content)
	definitions, err := h.Definition(uri, 17, 39)

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
	if d.Range.Start.Line != 6 {
		t.Error("expected definition start line to be 6, but was", d.Range.Start.Line)
	}
	if d.Range.Start.Character != 0 {
		t.Error("expected definition start character to be 0, but was", d.Range.Start.Character)
	}
	if d.Range.End.Line != 6 {
		t.Error("expected definition end line to be 6, but was", d.Range.End.Line)
	}
	if d.Range.End.Character != 0 {
		t.Error("expected definition end character to be 0, but was", d.Range.End.Character)
	}
}

func TestDefinition_AfterDocChange(t *testing.T) {
	h := newTestServer(t)
	uri, content := loadTestdata(t, "definition/internal/create.yaml")
	_, changedContent := loadTestdata(t, "definition/internal/create-changed.yaml")

	h.DidOpen(uri, "yaml", content)
	h.DidChange(uri, 2, changedContent)
	definitions, err := h.Definition(uri, 10, 34)

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
	if d.Range.Start.Line != 9 {
		t.Error("expected definition start line to be 6, but was", d.Range.Start.Line)
	}
	if d.Range.Start.Character != 0 {
		t.Error("expected definition start character to be 0, but was", d.Range.Start.Character)
	}
	if d.Range.End.Line != 9 {
		t.Error("expected definition end line to be 6, but was", d.Range.End.Line)
	}
	if d.Range.End.Character != 0 {
		t.Error("expected definition end character to be 0, but was", d.Range.End.Character)
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
