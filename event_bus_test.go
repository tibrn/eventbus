package eventbus

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	if got := New(&Options{Concurrency: 5}); got == nil {
		t.Errorf("New() returned nil")
	} else if len(got.list.queues) != 5 {
		t.Errorf("New() concurrency options is not applied")
	}
}

func TestEventBus_removeEvent(t *testing.T) {
	b := New()

	const (
		testEvent EventName = "test"
	)

	isCalled := false

	b.On(testEvent, func(payload []byte) {
		isCalled = true
	})

	b.removeEvent(testEvent)

	b.Emit(testEvent, []byte{})

	if isCalled {
		t.Error("EventBus.removeEvent() dosen't remove event from bus")
	}
}

func TestEventBus_Emit_On(t *testing.T) {
	b := New()

	const (
		testEvent               EventName = "test"
		testEventVaridiacParams           = "test2"
	)

	isCalled := false

	b.On(testEvent, func(payload []byte) {
		isCalled = true
	})

	waitCallback := b.Emit(testEvent, []byte{})

	err := waitCallback()

	if err != nil {

		t.Errorf("waitCallback got error %v", err)
	}

	if !isCalled {
		t.Error("EventBus.On()/.Emit() callbacks are not runned at all/ waitingCallback dosen't stop gorutine until all callbacks are called")
	}

	waitCallback = b.Emit(testEvent, []string{})

	err = waitCallback()

	if err == nil {

		t.Errorf("waitCallback should return error when payload is not correct")
	}

	_, err = b.On(testEvent, func(payload []string) {

	})

	if err == nil {

		t.Errorf("EventBus.On() should return error when signature dosen't match")
	}

	isCalled = false

	b.On(testEventVaridiacParams, func(a ...int) {
		isCalled = true
	})

	waitCallback = b.Emit(testEventVaridiacParams, 5, 6, 7)

	err = waitCallback()

	if err != nil {

		t.Errorf("waitCallback got error %v", err)
	}

	if !isCalled {
		t.Error("EventBus.On()/.Emit() callbacks are not runned at all when callback has variadic params/ waitingCallback dosen't stop gorutine until all callbacks are called")
	}

}

func Test_callbackRunner(t *testing.T) {
	if got := callbackRunner(&sync.WaitGroup{}, nil, nil, nil); got == nil {
		t.Errorf("callbackRunner() return shouldn't be nil")
	}
}

func Test_validateCallback(t *testing.T) {
	type args struct {
		task interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Accept functions",
			args: args{
				task: func() {},
			},
		},
		{
			name: "Don't accept struct",
			args: args{
				task: struct{}{},
			},
			wantErr: true,
		},
		{
			name: "Don't accept base type",
			args: args{
				task: 0,
			},
			wantErr: true,
		},
		{
			name: "Accept variadic function",
			args: args{
				task: func(...interface{}) {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateCallback(tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
