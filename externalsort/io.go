// +build !change

package externalsort

// LineReader ...
type LineReader interface {
	ReadLine() (string, error)
}

// LineWriter ...
type LineWriter interface {
	Write(l string) error
}
