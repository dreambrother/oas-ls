package operations

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreambrother/oas-ls/internal/textpos"
	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"gopkg.in/yaml.v3"
)

func Definition(ctx context.Context, log *slog.Logger, doc *document.Document, params *lsp.DefinitionParams) ([]lsp.Location, error) {
	line, ok := doc.Line(params.Position.Line)
	if !ok {
		log.WarnContext(ctx, "line not found", "line", params.Position.Line)
		return nil, nil
	}
	ref, ok := parseRef(line)
	if !ok {
		log.DebugContext(ctx, "line is not parseable", "line", line)
		return nil, nil
	}
	path, _ := strings.CutPrefix(string(doc.URI()), "file://")
	loc, ok := findLocation(ref, path)
	if !ok {
		log.DebugContext(ctx, "definition is not found", "definition", ref, "source", path)
		return nil, nil
	}
	return []lsp.Location{loc}, nil
}

func parseRef(line string) (string, bool) {
	refNode := &yaml.Node{}
	err := yaml.Unmarshal([]byte(line), refNode)
	if err != nil {
		return "", false
	}
	if refNode.Kind == yaml.DocumentNode {
		refNode = refNode.Content[0]
		if refNode.Kind == yaml.MappingNode && refNode.Content[0].Value == "$ref" {
			return refNode.Content[1].Value, true
		}
	}
	return "", false
}

func findLocation(ref string, docPath string) (lsp.Location, bool) {
	path, specType, found := strings.Cut(ref, "#/")
	if !found {
		return lsp.Location{}, false
	}
	// TODO check store

	targetPath := filepath.Join(filepath.Dir(docPath), filepath.Clean(path))
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return lsp.Location{}, false
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return lsp.Location{}, false
	}

	n, ok := resolveType(&doc, specType)
	if !ok {
		return lsp.Location{}, false
	}

	pos := position(data, n)
	return lsp.Location{
		URI: lsp.DocumentURI("file://" + targetPath),
		Range: lsp.Range{
			Start: pos,
			End:   pos,
		},
	}, true
}

func resolveType(n *yaml.Node, t string) (*yaml.Node, bool) {
	if n.Kind == yaml.DocumentNode {
		if len(n.Content) == 0 {
			return nil, false
		}
		n = n.Content[0]
	}
	if n.Kind != yaml.MappingNode {
		return nil, false
	}
	for i := 0; i < len(n.Content); i += 2 {
		if n.Content[i].Value == t {
			return n.Content[i], true
		}
	}
	return nil, false
}

// Position converts a yaml.v3 node Position into an LSP Position. yaml.v3
// reports 1-based lines and code-point columns, while LSP (given the advertised
// UTF-16 encoding) uses 0-based lines and UTF-16 code-unit offsets.
func position(data []byte, n *yaml.Node) lsp.Position {
	line := n.Line - 1
	var text string
	if lines := strings.Split(string(data), "\n"); line >= 0 && line < len(lines) {
		text = lines[line]
	}
	return lsp.Position{Line: line, Character: textpos.UTF16Column(text, n.Column-1)}
}
