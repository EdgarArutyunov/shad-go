// +build !solution

package keylock

import "fmt"

// KeyLock ...
type KeyLock struct {
	mt chan struct{}
	st map[string]chan struct{}
}

// New ...
func New() *KeyLock {
	res := &KeyLock{
		mt: make(chan struct{}, 1),
		st: make(map[string]chan struct{}),
	}
	res.mt <- struct{}{}
	return res
}

var debug = false

func print(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// LockKeys ...
func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	for {
		select {
		case <-l.mt:

		case <-cancel:
			return true, nil
		}
		stop := -1

	Loop:
		for i, key := range keys {
			_, ok := l.st[key]
			if !ok {
				l.st[key] = make(chan struct{})
				continue
			}

			select {
			case <-l.st[keys[i]]:
				l.st[key] = make(chan struct{})
				continue
			default:
				stop = i
				break Loop
			}
		}
		if stop != -1 {
			for i := 0; i != stop; i++ {
				close(l.st[keys[i]])
			}

			l.mt <- struct{}{}

			select {
			case <-l.st[keys[stop]]:
				continue
			case <-cancel:
				return true, nil
			}
		}

		l.mt <- struct{}{}

		return false, func() {
			for _, key := range keys {
				close(l.st[key])
			}
		}
	}
}
