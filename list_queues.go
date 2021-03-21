package eventbus

import (
	"sync"
)

type ListQueues struct {
	*Semaphore

	queues []*QueueCallbacks
	len    int

	mtx *sync.Mutex
}

func NewListQueues(nrQueues int, limitRunner int) *ListQueues {

	if nrQueues <= 0 {
		panic("NewListQueues(): number of queues must be greater than 0")
	}

	queues := make([]*QueueCallbacks, nrQueues)
	semaphore := NewSemaphore(nrQueues)

	for i := range queues {
		queues[i] = NewQueue(limitRunner)
	}

	return &ListQueues{
		queues:    queues,
		Semaphore: semaphore,
		mtx:       &sync.Mutex{},
		len:       nrQueues,
	}
}

func (c *ListQueues) Push(x *QueueCallbacks) {
	c.mtx.Lock()

	c.len++
	for i := 0; i < c.len; i++ {

		if c.queues[i] == nil {
			c.queues[i] = x
			break
		} else if x.WaitingCallbacks() < c.queues[i].WaitingCallbacks() {

			for j := i; j < c.len; j++ {
				if c.queues[j] == nil {
					c.queues[j] = x
					break
				} else {
					c.queues[i], x = x, c.queues[i]
				}
			}

			break
		}
	}

	c.mtx.Unlock()
}

func (c *ListQueues) Pop() *QueueCallbacks {

	var (
		x *QueueCallbacks
	)

	c.mtx.Lock()

	c.queues[0], x = x, c.queues[0]

	for i := 1; i < c.len; i++ {
		c.queues[i-1], c.queues[i] = c.queues[i], c.queues[i-1]
	}
	c.len--

	c.mtx.Unlock()

	return x
}
