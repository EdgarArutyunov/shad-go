// +build !solution

package otp

import (
	"io"
)

// MyReader ...
type MyReader struct {
	r io.Reader
	g io.Reader
	i int
}

// MyWriter ...
type MyWriter struct {
	w io.Writer
	g io.Reader
}

// NewReader ...
func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &MyReader{
		r: r,
		g: prng,
	}
}

// NewWriter ...
func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &MyWriter{
		w: w,
		g: prng,
	}
}

// [> Methods <]

// Read for MyReader ...
func (r *MyReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)

	buf := make([]byte, n)
	_, _ = r.g.Read(buf)

	for i := 0; i < n; i++ {
		p[i] ^= buf[i]
	}
	return n, err
}

// Write for MyWriter ...
func (w *MyWriter) Write(p []byte) (n int, err error) {
	ln := len(p)
	buf := make([]byte, ln)
	_, _ = w.g.Read(buf)

	for i := 0; i < ln; i++ {
		buf[i] ^= p[i]
	}
	return w.w.Write(buf)
}
