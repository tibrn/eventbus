package eventbus

import (
	"testing"
	"time"
)

func TestNewRunnerFn(t *testing.T) {

	t.Run("Not nil", func(t *testing.T) {
		if got := NewRunnerFn(time.Second, 1); got == nil {
			t.Errorf("NewRunnerFn() got nil")
		}
	})
}

func TestPorcessFn_Stop(t *testing.T) {

	maxTime := time.Millisecond * 10
	pf := NewRunnerFn(maxTime, 5)

	for i := 0; i < 5; i++ {
		pf.Next(func() {
			time.Sleep(maxTime * 2)
		})
	}
	pf.Stop()

	time.Sleep(maxTime * 4)

	if pf.currentNrFnWorkers != 0 {
		t.Errorf("Runner.Stop() current number of workers should be 0 but is %d", pf.currentNrFnWorkers)
	}

}

func TestPorcessFn_Next(t *testing.T) {

	isFinished := false

	maxTime := time.Millisecond * 1
	limit := 3
	pf := NewRunnerFn(maxTime, limit)

	t1 := time.Now()

	pf.Next(func() {
		time.Sleep(maxTime * 2)
	})

	elapsed := time.Since(t1)

	if elapsed > maxTime+time.Millisecond {
		t.Errorf("Runner.Next() dosen't free blocking after maxTime")
	}

	if pf.currentNrFnWorkers < 0 {
		t.Errorf("Runner.Next() dosen't count each worker")
	}

	time.Sleep(maxTime * 2)

	if pf.currentNrFnWorkers < 0 {
		t.Errorf("Runner.Next() releases last worker to early")
	}

	t1 = time.Now()
	go func() {
		time.Sleep(maxTime / 2)
		if pf.currentNrFnWorkers > 1 {
			t.Errorf("Runner.Next() increases workers unecessary")
		}
	}()

	pf.Next(func() {
		time.Sleep(maxTime)
	})

	time.Sleep(maxTime)

	go func() {
		for !isFinished {
			time.Sleep(time.Nanosecond)
			if pf.currentNrFnWorkers > int32(limit) {
				t.Errorf("Runner.Next() increases workers over the limit")
				return
			}
		}
	}()

	for i := 0; i < limit+1; i++ {

		if i == limit {
			t1 = time.Now()
		}
		pf.Next(func() {
			time.Sleep(maxTime * time.Duration(limit) * 5)
		})
		if i == limit {
			elapsed = time.Since(t1)
		}

	}

	if elapsed < maxTime*2 {
		t.Errorf("Runner.Next() dosen't block when limit is reached")
	}

	isFinished = true
}
