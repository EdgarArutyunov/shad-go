// +build !solution

package spacecollapse

import (
	"bytes"
	"unicode"
)

// CollapseSpaces ...
func CollapseSpaces(input string) string {
	var buf bytes.Buffer

	lastSpace := false
	for _, r := range input {
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
