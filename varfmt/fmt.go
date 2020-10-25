// +build !solution

package varfmt

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Sprintf ...
func Sprintf(format string, args ...interface{}) string {
	var buf, curArg bytes.Buffer

	const Nul = 0
	const Open = 1

	st := Nul
	argNum := 0
	for i := 0; i < len(format); {
		r, sz := utf8.DecodeRuneInString(format[i:])
		i += sz
		switch st {
		case Nul:
			if r == '{' {
				st = Open
			} else {
				buf.WriteRune(r)
			}
		case Open:
			if r == '}' {
				if len(curArg.String()) != 0 {
					arg, _ := strconv.Atoi(curArg.String())
					curArg.Reset()
					buf.WriteString(fmt.Sprint(args[arg]))
				} else {
					buf.WriteString(fmt.Sprint(args[argNum]))
				}
				argNum++
				st = Nul
			} else {
				curArg.WriteRune(r)
			}
		}
	}
	return buf.String()
}
