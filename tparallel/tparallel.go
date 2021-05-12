// +build !solution

package tparallel

import "sync"

// T ...
type T struct {
	waitParent     chan struct{} // wa
	parallel       chan struct{} // send to parent that we are in parallel waiting
	done           chan struct{} // test-function and childs are closed
	functionClosed chan struct{} // close this ch when test-function closed
	//signal for childs in parallel waiting

	wg *sync.WaitGroup // wait childs
}

// Parallel ...
func (t *T) Parallel() {
	t.parallel <- struct{}{}
	<-t.waitParent
}

// Run ...
func (t *T) Run(subtest func(t *T)) {
	// когда t завершится, - st можно стартовать

	st := &T{
		parallel:       make(chan struct{}),
		done:           make(chan struct{}),
		functionClosed: make(chan struct{}),
		waitParent:     t.functionClosed,
		wg:             &sync.WaitGroup{},
	}

	t.wg.Add(1)

	go func(st, t *T, subtest func(t *T)) {
		subtest(st)
		close(st.functionClosed) // send signal to childs
		st.wg.Wait()             // Wait childs

		defer close(st.done) // send signal to run
		defer t.wg.Done()    // end in done
	}(st, t, subtest)

	select {
	case <-st.parallel:
	case <-st.done:
	}
}

// Run ...
func Run(topTests []func(t *T)) {
	runParallelTests := make(chan struct{})
	wg := &sync.WaitGroup{}

	for _, tst := range topTests {
		t := &T{
			parallel:       nil,
			functionClosed: runParallelTests, // брейн фак хук в печень
			// для обычных тестов нам не нужны
			done:       nil,
			waitParent: nil,
			wg:         wg,
		}

		t.Run(tst)
	}

	close(runParallelTests)
	wg.Wait()
}
