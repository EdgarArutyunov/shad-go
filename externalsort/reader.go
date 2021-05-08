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
	bytes     int // read bytes
	readBytes int // how much we already read to bytes.Buffer in Readline
	err       error
}

// ReadLine ...
func (mr *MyLineReader) ReadLine() (string, error) {
	var buf bytes.Buffer
	for {
		if mr.bytes == mr.readBytes {
			if mr.err != nil {
				return buf.String(), mr.err
			}
			mr.bytes, mr.err = mr.r.Read(mr.buf)
			mr.readBytes = 0
			continue
		}
		newLinePos := bytes.IndexAny(mr.buf[mr.readBytes:mr.bytes], "\n")
		if newLinePos == -1 {
			buf.Write(mr.buf[mr.readBytes:mr.bytes])
			mr.readBytes = mr.bytes
		} else {
			buf.Write(mr.buf[mr.readBytes : mr.readBytes+newLinePos]) // [ : \n)
			mr.readBytes += newLinePos + 1                            // (\n : ]
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
		bytes:     0,
		err:       nil,
	}
}
