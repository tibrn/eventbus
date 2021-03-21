package eventbus

import (
	"sync"
	"sync/atomic"
	"time"
)

type PorcessFn struct {
	currentState State

	chFn               chan FnCallback
	currentNrFnWorkers int32
	semaphore          *Semaphore

	limitTicker *time.Ticker
	ticker      *time.Ticker

	group *sync.WaitGroup

	mtx *sync.Mutex

	maxTime  time.Duration
	isStoped bool
}

func NewRunnerFn(maxTime time.Duration, limit int) *PorcessFn {

	return &PorcessFn{
		ticker:      time.NewTicker(time.Nanosecond * 1),
		limitTicker: time.NewTicker(maxTime),
		maxTime:     maxTime,
		mtx:         &sync.Mutex{},
		group:       &sync.WaitGroup{},
		chFn:        make(chan FnCallback),
		semaphore:   NewSemaphore(limit),
	}
}

func (pf *PorcessFn) Stop() WaitCallback {
	pf.isStoped = true

	return func() error {
		pf.group.Wait()
		return nil
	}
}

func (pf *PorcessFn) acquire() {
	pf.group.Add(1)
	pf.semaphore.Acquire()
	atomic.AddInt32(&pf.currentNrFnWorkers, 1)
}

func (pf *PorcessFn) release() {
	atomic.AddInt32(&pf.currentNrFnWorkers, -1)
	pf.semaphore.Release()
	pf.group.Done()
}

func (pf *PorcessFn) Next(fn FnCallback) error {

	if fn == nil {
		return ErrInvalidFnCallback
	}

	pf.mtx.Lock()
	defer func() {
		pf.limitTicker.Stop()
		pf.ticker.Stop()
		pf.mtx.Unlock()
	}()

	//Don't accept any new func to PorcessFn if one is already running
	if pf.currentState.IsStilRunning() {
		return ErrStillRunning
	}

	//Reset tickers
	pf.limitTicker.Reset(pf.maxTime)
	pf.ticker.Reset(time.Microsecond * 10)

	//Set new state
	pf.currentState = RUNNING
	isDone := false

	if atomic.LoadInt32(&pf.currentNrFnWorkers) == 0 || pf.currentState.IsOverMaxTime() {

		//Block any new FnCallack if limit is reached
		pf.acquire()
		go func() {
			var (
				fn     FnCallback
				ticker = time.NewTicker(time.Millisecond * 10)
			)

			defer func() {
				ticker.Stop()
				pf.release()
			}()

			for {

				select {
				case fn = <-pf.chFn:
					isDone = false
					fn()
					isDone = true
					if atomic.LoadInt32(&pf.currentNrFnWorkers) == 1 {
						//Set current state done when is only one worker to be able to accept new FnCallback
						pf.currentState = DONE
					} else if atomic.LoadInt32(&pf.currentNrFnWorkers) > 1 {
						return
					}
					break
				case <-ticker.C:
					if pf.isStoped || atomic.LoadInt32(&pf.currentNrFnWorkers) > 1 {
						return
					}
				}

			}
		}()
	}

	pf.chFn <- fn

	//Check for done/overTime
	for {
		select {
		case <-pf.limitTicker.C:
			pf.currentState = OVERMAXTIME
			return nil

		case <-pf.ticker.C:
			if isDone {
				pf.currentState = DONE
				return nil
			}
		}
	}
}
