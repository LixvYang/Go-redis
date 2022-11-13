package wait

import (
	"sync"
	"time"
)

type Wait struct {
	wg sync.WaitGroup
}

func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

func (w *Wait) Done() {
	w.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
func (w *Wait) Wait() {
	w.wg.Wait()
}

func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{}, 1)

	go func() {
		defer close(c)
		w.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return false
	case <- time.After(timeout):
		return true
	}
}
