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
		again := false
	Loop:
		for i, key := range keys {
			_, ok := l.st[key]
			if !ok {
				l.st[key] = make(chan struct{}, 1)
				continue
			}

			select {
			case <-l.st[key]:
			default:
				for j := 0; j < i; j++ {
					l.st[keys[j]] <- struct{}{}
				}
				l.mt <- struct{}{}
				again = true
				break Loop
			}
		}
		if again {
			continue
		}

		l.mt <- struct{}{}

		return false, func() {
			for _, key := range keys {
				l.st[key] <- struct{}{}
			}
		}
	}
}
