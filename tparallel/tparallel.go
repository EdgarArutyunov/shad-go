// +build !solution

package tparallel

// T ...
type T struct {
	chGo                chan struct{}
	waitTokenFromParent chan struct{}
	sendTokensToChild   []chan struct{}
	waitEndOfChilds     []chan struct{}
	sendEndToParent     chan struct{}
}

// Parallel ...
func (t *T) Parallel() {
	t.chGo <- struct{}{}
	<-t.waitTokenFromParent
}

// Run ...
func (t *T) Run(subtest func(t *T)) {

	curEndID := len(t.waitEndOfChilds)
	t.waitEndOfChilds = append(
		t.waitEndOfChilds,
		make(chan struct{}, 1),
	)

	t.sendTokensToChild = append(
		t.sendTokensToChild,
		make(chan struct{}, 1),
	)

	subT := &T{
		chGo:                make(chan struct{}, 1),
		waitTokenFromParent: t.sendTokensToChild[curEndID],
		sendTokensToChild:   make([]chan struct{}, 0),
		waitEndOfChilds:     make([]chan struct{}, 0),
		sendEndToParent:     t.waitEndOfChilds[curEndID],
	}

	go func(subT *T, t *T) {
		subtest(subT)
		defer func() {
			for _, send := range subT.sendTokensToChild {
				send <- struct{}{}
			}

			for _, getEnd := range subT.waitEndOfChilds {
				<-getEnd
			}

			subT.sendEndToParent <- struct{}{}
		}()
	}(subT, t)

	select {
	case <-t.waitEndOfChilds[curEndID]:
		close(t.waitEndOfChilds[curEndID])

	case <-subT.chGo:
	}
}

// Run ...
func Run(topTests []func(t *T)) {
	t := &T{
		chGo:                nil,
		waitTokenFromParent: nil,
		sendTokensToChild:   make([]chan struct{}, 0),
		waitEndOfChilds:     make([]chan struct{}, 0),
		sendEndToParent:     nil,
	}

	for _, tst := range topTests {
		t.Run(tst)
	}

	for _, send := range t.sendTokensToChild {
		send <- struct{}{}
	}

	for _, wait := range t.waitEndOfChilds {
		<-wait
	}
}
