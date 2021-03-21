package eventbus

import (
	"sync"
	"sync/atomic"
	"time"
)

type FnCallback func()

type QueueCallbacks struct {
	isRunning        bool
	startOnce        sync.Once
	waitingCallbacks uint32
	queue            []FnCallback
	mtx              *sync.Mutex

	runnerFn    *PorcessFn
	runnerLimit int
}

func NewQueue(runnerLimit int) *QueueCallbacks {
	q := &QueueCallbacks{
		startOnce:        sync.Once{},
		queue:            make([]FnCallback, runnerLimit),
		mtx:              &sync.Mutex{},
		runnerFn:         NewRunnerFn(time.Millisecond*10, runnerLimit),
		runnerLimit:      runnerLimit,
		waitingCallbacks: 0,
	}

	return q
}

func (q *QueueCallbacks) Set(fn FnCallback) {

	atomic.AddUint32(&q.waitingCallbacks, 1)
	q.fitOrExtand(fn)

	q.start()

}

func (q *QueueCallbacks) fitOrExtand(fn FnCallback) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	for i := len(q.queue) - 1; i >= 0; i-- {
		if q.queue[i] == nil {
			q.queue[i] = fn
			return
		}
	}

	q.queue = append(q.queue, fn)

}

func (q *QueueCallbacks) WaitingCallbacks() uint32 {
	return atomic.LoadUint32(&q.waitingCallbacks)
}

const (
	timeTicker = 1 * time.Second
)

func (q *QueueCallbacks) start() {

	if q.isRunning {
		return
	}

	q.isRunning = true

	q.startOnce.Do(func() {

		//Consumer
		go func() {

			defer func() {
				if err := recover(); err != nil {

				}
				q.resetQueue()
			}()

			var (
				fn           FnCallback
				i            int
				lenQ         int
				timeEndQueue *time.Time
				loopCheck    int
			)

			for q.isRunning {
				i = 0
				lenQ = len(q.queue)
				for ; i < lenQ; i++ {

					if q.queue[i] == nil {
						continue
					}

					if timeEndQueue != nil {
						timeEndQueue = nil
					}

					fn, q.queue[i] = q.queue[i], fn

					q.runnerFn.Next(fn)

					fn = nil

				}

				if timeEndQueue == nil {
					t := time.Now()
					timeEndQueue = &t
					loopCheck = 0
				} else if loopCheck > 6 && timeEndQueue.After(time.Now().Add(-1*timeTicker)) {
					return
				}
				loopCheck++
				time.Sleep(time.Nanosecond * 50)
			}

		}()
	})

}

func (q *QueueCallbacks) resetQueue() {

	q.queue = make([]FnCallback, q.runnerLimit)

	q.startOnce = sync.Once{}

	q.isRunning = false

	atomic.SwapUint32(&q.waitingCallbacks, 0)

}
