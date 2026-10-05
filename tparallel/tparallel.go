//go:build !solution

package tparallel

import (
	"sync"
)

type state int

const (
	running  state = iota
	parallel       // called Parallel(), suspended
	done           // fully finished including parallel children
)

type T struct {
	parent    *T
	cond      *sync.Cond
	state     state
	parallels []*T
	parWg     sync.WaitGroup
	release   chan struct{}
}

func newT(parent *T) *T {
	t := &T{
		parent:  parent,
		release: make(chan struct{}),
	}
	t.cond = sync.NewCond(&sync.Mutex{})
	return t
}

func updateState(t *T, s state) {
	t.cond.L.Lock()
	t.state = s
	t.cond.Signal()
	t.cond.L.Unlock()
}

func (t *T) Parallel() {
	t.parent.cond.L.Lock()
	t.parent.parallels = append(t.parent.parallels, t)
	t.parent.parWg.Add(1)
	t.parent.cond.L.Unlock()

	updateState(t, parallel)
	<-t.release
}

func (t *T) Run(subtest func(*T)) {
	subT := newT(t)

	go func() {
		subtest(subT)

		// unblock subT's parallel children
		subT.cond.L.Lock()
		parallels := subT.parallels
		subT.cond.L.Unlock()

		for _, child := range parallels {
			close(child.release)
		}
		subT.parWg.Wait()
		
		prevState := subT.state
		updateState(subT, done)

		if prevState == parallel {
			t.parWg.Done()
		}
	}()

	subT.cond.L.Lock()
	for subT.state == running {
		subT.cond.Wait()
	}
	st := subT.state
	subT.cond.L.Unlock()

	// parallel: return 
	if st == parallel {
		return
	}
	
	// sequential: wait for test state to be `done`
	subT.cond.L.Lock()
	for subT.state != done {
		subT.cond.Wait()
	}
	subT.cond.L.Unlock()
}

func Run(topTests []func(*T)) {
	t := newT(nil)
	t.parent = t

	for _, tt := range topTests {
		t.Run(tt)
	}

	t.cond.L.Lock()
	parallels := t.parallels
	t.cond.L.Unlock()
	for _, child := range parallels {
		close(child.release)
	}
	t.parWg.Wait()
}
