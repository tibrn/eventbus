package eventbus

import (
	"reflect"
	"sync"
)

//Remove callback from listener
type RemoveCallback func()

//Wait for all callbacks to run with payload
type WaitCallback func() error

//Bus
type EventBus struct {
	events map[EventName]*Listeners
	mtx    *sync.RWMutex

	list *ListQueues
}

//Create new EventBus
func New(options ...*Options) *EventBus {

	finalOptions := defaultOrOptions(options...)

	return &EventBus{
		events: make(map[EventName]*Listeners),
		mtx:    &sync.RWMutex{},
		list:   NewListQueues(finalOptions.Concurrency, finalOptions.RunnerConcurrency),
	}
}

func (b *EventBus) removeEvent(event EventName) {
	b.mtx.Lock()
	delete(b.events, event)

	//Set old listeners for GC
	if len(b.events) == 0 {
		b.events = make(map[EventName]*Listeners)
	}

	b.mtx.Unlock()
}

//On stores new listener for event X
func (b *EventBus) On(event EventName, fn interface{}) (RemoveCallback, error) {

	fnValue, err := validateCallback(fn)

	if err != nil {
		return nil, err
	}

	{
		var listeners *Listeners

		b.mtx.RLock()

		if _, ok := b.events[event]; !ok {

			b.events[event] = NewListeners(func() {
				b.removeEvent(event)
			})
		}

		listeners = b.events[event]

		b.mtx.RUnlock()

		if !listeners.IsSame(fnValue) {
			return nil, ErrSignature
		}

		listeners.Add(fnValue)

	}

	return func() {
		b.mtx.RLock()
		listeners := b.events[event]
		b.mtx.RUnlock()

		if listeners == nil {
			return
		}

		listeners.Delete(fnValue)
	}, nil
}

//EmitEvent... emit event to all bus callbacks
func (b *EventBus) EmitEvent(event Event) WaitCallback {
	return b.Emit(event.Name, event.Payload)
}

func defaultWaitCallback() error {
	return nil
}

//Emit... emits payload to all listeners for event X
func (b *EventBus) Emit(event EventName, payload ...interface{}) WaitCallback {
	b.mtx.RLock()

	var listeners *Listeners
	if list, ok := b.events[event]; !ok {
		return defaultWaitCallback
	} else {
		listeners = list
	}
	b.mtx.RUnlock()

	group := sync.WaitGroup{}

	group.Add(1)

	var err error

	go func() {

		b.list.Acquire()
		defer b.list.Release()

		queue := b.list.Pop()

		defer b.list.Push(queue)
		defer group.Done()

		callbacks := listeners.Callbacks()

		lenCallbacks := len(callbacks)

		if lenCallbacks < 1 {
			return
		}

		payloadTransformed, isValid := listeners.IsValidPayload(payload)

		if !isValid {
			err = ErrInvalidPayload
			return
		}

		group.Add(lenCallbacks)

		for _, callback := range callbacks {

			queue.Set(callbackRunner(&group, listeners, callback, payloadTransformed))
		}

	}()

	return func() error {
		group.Wait()

		return err
	}

}

func callbackRunner(group *sync.WaitGroup, listeners *Listeners, callback *reflect.Value, payload []reflect.Value) func() {

	return func() {
		defer func() {
			if err := recover(); err != nil {
				listeners.Delete(callback)
			}
			group.Done()
		}()

		callback.Call(payload)

	}
}

func validateCallback(fn interface{}) (*reflect.Value, error) {
	v := reflect.ValueOf(fn)
	t := v.Type()

	// Fn must be a function
	if t.Kind() != reflect.Func {
		return &v, ErrFnMustBeFunc
	}

	return &v, nil
}
