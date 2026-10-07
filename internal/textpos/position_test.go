package textpos

import "testing"

func TestUTF16Column(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		codePoints int
		want       int
	}{
		{
			name:       "empty line",
			line:       "",
			codePoints: 3,
			want:       0,
		},
		{
			name:       "ascii: each rune is one code unit",
			line:       "hello",
			codePoints: 3,
			want:       3,
		},
		{
			name:       "ascii: zero code points",
			line:       "hello",
			codePoints: 0,
			want:       0,
		},
		{
			name:       "ascii: code points exceed the line length",
			line:       "hello",
			codePoints: 100,
			want:       5,
		},
		{
			name:       "ascii: negative code points",
			line:       "hello",
			codePoints: -1,
			want:       0,
		},
		{
			name:       "cyrillic: BMP runes are one code unit each",
			line:       "привет",
			codePoints: 4,
			want:       4,
		},
		{
			name:       "cyrillic: whole line",
			line:       "привет",
			codePoints: 6,
			want:       6,
		},
		{
			name:       "cjk: BMP ideographs are one code unit each",
			line:       "日本語",
			codePoints: 2,
			want:       2,
		},
		{
			name:       "cjk: whole line",
			line:       "日本語",
			codePoints: 3,
			want:       3,
		},
		{
			name:       "emoji: non-BMP rune is a surrogate pair",
			line:       "😀",
			codePoints: 1,
			want:       2,
		},
		{
			name:       "emoji: several non-BMP runes",
			line:       "😀😀",
			codePoints: 2,
			want:       4,
		},
		{
			name:       "emoji: BMP symbol stays one code unit",
			line:       "☺",
			codePoints: 1,
			want:       1,
		},
		{
			name:       "emoji: base rune plus variation selector",
			line:       "❤️",
			codePoints: 2,
			want:       2,
		},
		{
			name:       "cjk extension b: non-BMP ideograph",
			line:       "𠀀",
			codePoints: 1,
			want:       2,
		},
		{
			name:       "mixed: ascii, cyrillic, emoji, cjk up to the emoji",
			line:       "aп😀日",
			codePoints: 2,
			want:       2,
		},
		{
			name:       "mixed: ascii, cyrillic, emoji, cjk including the emoji",
			line:       "aп😀日",
			codePoints: 3,
			want:       4,
		},
		{
			name:       "mixed: entire line",
			line:       "aп😀日",
			codePoints: 4,
			want:       5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UTF16Column(tt.line, tt.codePoints); got != tt.want {
				t.Errorf("UTF16Column(%q, %d) = %d, want %d", tt.line, tt.codePoints, got, tt.want)
			}
		})
	}
}
