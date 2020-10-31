package externalsort

import (
	"fmt"
	"io"
)

// MyLineWriter ...
type MyLineWriter struct {
	w io.Writer
}

// Write ...
func (mw *MyLineWriter) Write(l string) error {
	buf := []byte(l)
	buf = append(buf, '\n')
	n, err := mw.w.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("not all data")
	}
	return nil
}

// NewWriter ...
func NewWriter(w io.Writer) LineWriter {
	return &MyLineWriter{
		w: w,
	}
}
