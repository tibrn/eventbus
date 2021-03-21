package eventbus

import (
	"math"
	"testing"
	"time"
)

func TestSemaphore_Aquire_Release(t *testing.T) {
	s := NewSemaphore(1)

	s.Acquire()

	isDone := false
	go func() {
		time.Sleep(time.Millisecond * 50)
		s.Release()
		time.Sleep(time.Millisecond * 5)
		if !isDone {
			t.Error("Semaphore.Relase() dosen't free resources")
		}
	}()
	t1 := time.Now()
	s.Acquire()

	if time.Since(t1) < time.Millisecond*40 {
		t.Error("Semaphore.Acquire() dosen't block when limit of resouces is released")
	}

	isDone = true
}

func TestNewSemaphore(t *testing.T) {
	if got := NewSemaphore(-1); got == nil {
		t.Errorf("NewSemaphore() shouldn't return nil")
	}

	if got := NewSemaphore(math.MaxInt32); got == nil {
		t.Errorf("NewSemaphore() shouldn't return nil")
	}
}
