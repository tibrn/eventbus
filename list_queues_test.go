package eventbus

import (
	"sync"
	"testing"
	"time"
)

func TestNewListQueues(t *testing.T) {
	if got := NewListQueues(1, 1); got == nil {
		t.Errorf("NewListQueues() got nil")
	}
}

func TestListQueues_Push_Pop(t *testing.T) {

	const (
		nrQueues = 2
	)
	q := NewListQueues(nrQueues, 2)

	group := sync.WaitGroup{}
	group.Add(nrQueues)
	for i := 0; i < nrQueues; i++ {
		q.Acquire()

		go func() {
			item := q.Pop()
			time.Sleep(20 * time.Millisecond)
			q.Push(item)
			q.Release()
			group.Done()
		}()
	}

	time.Sleep(2 * time.Millisecond)
	t1 := time.Now()
	q.Acquire()

	item := q.Pop()

	if time.Since(t1) < time.Millisecond*10 {
		t.Error("ListQueues.Pop() dosen't limit the gorutines which cand have a queue in the same time")
	}

	q.Push(item)
	q.Release()

	group.Wait()

	for i := range q.queues {
		if q.queues[i] == nil {
			t.Error("ListQueues.Pop()/Push() dosen't work properly")
		}
	}
}
