package externalsort

import (
	"fmt"
	"io"
)

var (
	ErrNotAllDataWasSent = fmt.Errorf("not all data was sent")
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

	switch {
	case err != nil:
		return fmt.Errorf("error in writing: %w", err)

	case n != len(buf):
		return ErrNotAllDataWasSent

	default:
		return nil
	}
}

// NewWriter ...
func NewWriter(w io.Writer) LineWriter {
	return &MyLineWriter{
		w: w,
	}
}
