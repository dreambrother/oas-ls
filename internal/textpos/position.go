package textpos

import (
	"github.com/owenrumney/go-lsp/lsp"
)

const PositionEncoding = lsp.PositionEncodingUTF16

// utf16Column returns the number of UTF-16 code units occupied by the first
// codePoints runes of line.
func UTF16Column(line string, codePoints int) int {
	units := 0
	for _, r := range line {
		if codePoints <= 0 {
			break
		}
		codePoints--
		if r > 0xFFFF {
			units += 2
		} else {
			units++
		}
	}
	return units
}
