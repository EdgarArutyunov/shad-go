// +build !solution

package externalsort

import (
	"io"
	"os"
	"sort"
)

var debug = false

// Sort ...
func Sort(w io.Writer, in ...string) error {
	readers := make([]LineReader, 0)

	for _, fname := range in {
		fr, err := os.Open(fname)
		if err != nil {
			return err
		}

		mr := NewReader(fr)

		lines := make([]string, 0)

		for {
			line, err := mr.ReadLine()
			if err == nil {
				lines = append(lines, line)
			} else if err == io.EOF {
				if line != "" {
					lines = append(lines, line)
				}
				break
			} else {
				return err
			}
		}

		err = fr.Close()
		if err != nil {
			return err
		}

		sort.Slice(lines, func(i, j int) bool {
			return lines[i] < lines[j]
		})

		fw, err := os.Create(fname)
		if err != nil {
			return err
		}

		mw := NewWriter(fw)

		for _, line := range lines {
			err = mw.Write(line)
			if err != nil {
				return err
			}
		}
		err = fw.Close()
		if err != nil {
			return err
		}

		fr, err = os.Open(fname)
		defer fr.Close()
		readers = append(readers, NewReader(fr))
	}

	mw := NewWriter(w)
	return Merge(mw, readers...)
}
