package eventbus

import (
	"sync"
	"testing"
	"time"
)

func TestQueueCallbacks_start(t *testing.T) {
	queue := NewQueue(2)

	queue.start()

	isDone := false

	queue.Set(func() {
		isDone = true
	})

	time.Sleep(time.Millisecond * 1)

	if !isDone {
		t.Error("QueueCallbacks.start() dosen't start the consumer")
	}

}

func TestQueueCallbacks_WaitingCallbacks(t *testing.T) {
	type fields struct {
		isRunning        bool
		startOnce        sync.Once
		waitingCallbacks uint32
		queue            []FnCallback
		mtx              *sync.Mutex
		runnerFn         *PorcessFn
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Same as underline waintingCallbacks uint32",
			fields: fields{
				waitingCallbacks: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QueueCallbacks{
				isRunning:        tt.fields.isRunning,
				startOnce:        tt.fields.startOnce,
				waitingCallbacks: tt.fields.waitingCallbacks,
				queue:            tt.fields.queue,
				mtx:              tt.fields.mtx,
				runnerFn:         tt.fields.runnerFn,
			}
			if got := q.WaitingCallbacks(); got != tt.want {
				t.Errorf("QueueCallbacks.WaitingCallbacks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueueCallbacks_resetQueue(t *testing.T) {
	type fields struct {
		isRunning        bool
		startOnce        sync.Once
		waitingCallbacks uint32
		queue            []FnCallback
		mtx              *sync.Mutex
		runnerFn         *PorcessFn
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Dosen't panic",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QueueCallbacks{
				isRunning:        tt.fields.isRunning,
				startOnce:        tt.fields.startOnce,
				waitingCallbacks: tt.fields.waitingCallbacks,
				queue:            tt.fields.queue,
				mtx:              tt.fields.mtx,
				runnerFn:         tt.fields.runnerFn,
			}
			q.resetQueue()
		})
	}
}

func TestQueueCallbacks_fitOrExtand(t *testing.T) {
	q := NewQueue(2)

	lenQ := len(q.queue)

	q.fitOrExtand(func() {})

	if lenQ != len(q.queue) {
		t.Error("QueueCallbacks.fitOrExtand() extands queue when is not necessary")
	}

	for i := 0; i < lenQ; i++ {
		q.fitOrExtand(func() {})

	}

	if lenQ >= len(q.queue) {
		t.Error("QueueCallbacks.fitOrExtand() dosen't extands queue when is not necessary")
	}
}

func TestQueueCallbacks_Set(t *testing.T) {
	queue := NewQueue(2)

	queue.Set(func() {

	})

	if queue.WaitingCallbacks() != 1 {
		t.Error("QueueCallbacks.Set() dosen't increase the number of waitingCallbacks correctly")
	}

	for i := range queue.queue {
		if queue.queue[i] != nil {
			return
		}
	}

	t.Error("QueueCallbacks.Set() dosen't append the FnCallback to the underline slice of callbacks")
}

func TestNewQueue(t *testing.T) {
	if got := NewQueue(1); got == nil {
		t.Errorf("NewQueue() got nil")
	}
}
