// +build !solution

package spacecollapse

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

// CollapseSpaces ...
func CollapseSpaces(input string) string {
	var buf bytes.Buffer
	for i, lastSpace := 0, false; i < len(input); {
		r, sz := utf8.DecodeRuneInString(input[i:])
		i += sz

		if !unicode.IsSpace(r) {
			lastSpace = false
			buf.WriteRune(r)
		} else if !lastSpace {
			lastSpace = true
			buf.WriteRune(' ')
		}
	}
	return buf.String()
}
