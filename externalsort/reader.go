package externalsort

import (
	"bytes"
	"io"
)

const bufSize = 512

// MyLineReader ...
type MyLineReader struct {
	r         io.Reader
	buf       []byte
	readBytes int // readBytes 0, 1, 2 (how much)
	freeBytes int // 1, 2, 3
	err       error
}

// ReadLine ...
func (mr *MyLineReader) ReadLine() (string, error) {
	var buf bytes.Buffer
	for {
		if mr.readBytes == mr.freeBytes {
			if mr.err != nil {
				return buf.String(), mr.err
			}
			mr.freeBytes, mr.err = mr.r.Read(mr.buf)
			mr.readBytes = 0
			continue
		}
		firstEndLine := bytes.IndexAny(mr.buf[mr.readBytes:mr.freeBytes], "\n")
		if firstEndLine == -1 {
			end := mr.freeBytes
			buf.Write(mr.buf[mr.readBytes:end])
			mr.readBytes = end
		} else {
			end := mr.readBytes + firstEndLine
			buf.Write(mr.buf[mr.readBytes:end]) // ignore \n
			mr.readBytes = end + 1              // ignore \n
			return buf.String(), mr.err
		}
	}
}

// NewReader ...
func NewReader(r io.Reader) LineReader {
	return &MyLineReader{
		r:         r,
		buf:       make([]byte, bufSize),
		readBytes: 0,
		freeBytes: 0,
		err:       nil,
	}
}
