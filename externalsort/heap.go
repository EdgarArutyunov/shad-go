package externalsort

// ReadAndVal ...
type ReadAndVal struct {
	readerID int
	val      string
}

// An ReadAndValHeap is a min-heap of ints.
type ReadAndValHeap []ReadAndVal

func (h *ReadAndValHeap) Len() int           { return len(*h) }
func (h *ReadAndValHeap) Less(i, j int) bool { return (*h)[i].val < (*h)[j].val }
func (h *ReadAndValHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

// Push ...
func (h *ReadAndValHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(ReadAndVal))
}

// Pop ...
func (h *ReadAndValHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
