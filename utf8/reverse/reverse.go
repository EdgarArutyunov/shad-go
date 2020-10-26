// +build !solution

package reverse

// Reverse ...
func Reverse(input string) string {
	r := []rune(input)
	for i, j := 0, len(r)-1; i < j; i++ {
		r[i], r[j] = r[j], r[i]
		j--
	}
	return string(r)
}
