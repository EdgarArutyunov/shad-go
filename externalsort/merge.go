package externalsort

import (
	"container/heap"
	"io"
)

// Merge ...
func Merge(w LineWriter, readers ...LineReader) error {

	h := &ReadAndValHeap{}
	heap.Init(h)

	insertNewVal := func(readerID int) error {
		if readerID < 0 {
			return nil
		}
		val, err := readers[readerID].ReadLine()
		if err == nil {
			heap.Push(h, ReadAndVal{
				readerID: readerID,
				val:      val,
			})
			return nil
		}

		if err == io.EOF {
			if val != "" {
				heap.Push(h, ReadAndVal{
					readerID: -1,
					val:      val,
				})
			}
			return nil
		}

		return err
	}

	for i := range readers {
		if err := insertNewVal(i); err != nil {
			return err
		}
	}

	for h.Len() > 0 {
		readAndVal := heap.Pop(h)
		err := w.Write(readAndVal.(ReadAndVal).val)
		if err != nil {
			return err
		}
		if err := insertNewVal(readAndVal.(ReadAndVal).readerID); err != nil {
			return err
		}
	}
	return nil
}
