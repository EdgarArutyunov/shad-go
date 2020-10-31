// +build !solution

package lrucache

import "container/list"

// Lc ...
type Lc struct {
	cap int
	m   map[int]*list.Element
	l   *list.List
}

type pair struct {
	key, val int
}

// Get ...
func (lc *Lc) Get(key int) (int, bool) {
	elem, ok := lc.m[key]
	if !ok {
		return 0, false
	}
	lc.l.MoveToBack(elem)

	lc.m[key] = lc.l.Back()
	return lc.l.Back().Value.(pair).val, true
}

// Set ...
func (lc *Lc) Set(key, value int) {
	if lc.cap == 0 {
		return
	}

	elem, ok := lc.m[key]
	if ok {
		lc.l.MoveToBack(elem)
		lc.m[key] = lc.l.Back()
		lc.l.Back().Value = pair{
			key: key,
			val: value,
		}
		return
	}

	if len(lc.m) == lc.cap {
		delete(lc.m, lc.l.Front().Value.(pair).key)
		lc.l.Remove(lc.l.Front())
	}
	lc.m[key] = lc.l.PushBack(pair{
		key: key,
		val: value,
	})
}

// Range ...
func (lc *Lc) Range(f func(key, value int) bool) {
	for l := lc.l.Front(); l != nil; l = l.Next() {
		info := l.Value.(pair)
		if !f(info.key, info.val) {
			break
		}
	}
}

// Clear ...
func (lc *Lc) Clear() {
	lc.m = make(map[int]*list.Element)
	lc.l = list.New()
}

// New ...
func New(cap int) Cache {
	return &Lc{
		m:   make(map[int]*list.Element),
		l:   list.New(),
		cap: cap,
	}
}
