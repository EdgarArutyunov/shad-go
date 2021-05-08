package externalsort

import (
	"container/heap"
	"fmt"
	"io"
)

var (
	ErrHeapCastError = fmt.Errorf("heap contains unexpected value")
)

func closed_reader() int {
	return -1
}

// Merge ...
func Merge(w LineWriter, readers ...LineReader) error {

	h := &ReadAndValHeap{}
	heap.Init(h)

	insertNewValFromReader := func(readerID int) error {
		if readerID == closed_reader() {
			return nil
		}

		val, err := readers[readerID].ReadLine()

		switch err {
		case nil:
			heap.Push(h, ReadAndVal{
				readerID: readerID,
				val:      val,
			})
			return nil

		case io.EOF:
			if val != "" {
				heap.Push(h, ReadAndVal{
					readerID: -1,
					val:      val,
				})
			}
			return nil

		default:
			return err
		}
	}

	for i := range readers {
		if err := insertNewValFromReader(i); err != nil {
			return err
		}
	}

	for h.Len() > 0 {
		popValue, ok := heap.Pop(h).(ReadAndVal)

		if !ok {
			return ErrHeapCastError
		}

		if err := w.Write(popValue.val); err != nil {
			return err
		}

		if err := insertNewValFromReader(popValue.readerID); err != nil {
			return err
		}
	}

	return nil
}
