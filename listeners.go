package eventbus

import (
	"reflect"
	"sync"
)

type Listeners struct {
	callbacks map[*reflect.Value]*reflect.Value
	mtx       *sync.RWMutex

	remove RemoveCallback

	sliceCallbacks []*reflect.Value

	signature *reflect.Type
}

func isSameType(v1, v2 reflect.Type) bool {

	if v1.Kind() != v2.Kind() {
		return false
	}

	return v1.AssignableTo(v2)
}

func NewListeners(remove RemoveCallback) *Listeners {
	return &Listeners{
		callbacks: make(map[*reflect.Value]*reflect.Value),
		mtx:       &sync.RWMutex{},
		remove:    remove,
	}
}

func (l *Listeners) IsSame(v *reflect.Value) bool {
	l.mtx.Lock()

	if l.signature == nil {
		l.mtx.Unlock()
		return true
	}

	signature := *l.signature
	l.mtx.Unlock()

	vType := v.Type()

	lenIn := vType.NumIn()

	if signature.NumIn() != lenIn {
		return false
	}

	for i := 0; i < lenIn; i++ {
		if !isSameType(signature.In(i), vType.In(i)) {
			return false
		}
	}

	return true
}

func (l *Listeners) IsValidPayload(payload []interface{}) ([]reflect.Value, bool) {
	l.mtx.Lock()

	if l.signature == nil {
		l.mtx.Unlock()
		return nil, false
	}

	signature := *l.signature

	l.mtx.Unlock()

	var (
		isVariadic  = signature.IsVariadic()
		lastNrInput = signature.NumIn() - 1
		lenPayload  = len(payload)
		elem        reflect.Type
	)

	//Check for payload to have the minimum of required params
	if (!isVariadic && signature.NumIn() != lenPayload) ||
		(isVariadic && lenPayload < lastNrInput) {
		return nil, false
	}

	//Transform and validate payload
	transformedPayload := make([]reflect.Value, len(payload))

	for i, j := 0, 0; i < lenPayload; i++ {

		transformedPayload[i] = reflect.ValueOf(payload[i])

		if !isVariadic || j < lastNrInput {
			elem = signature.In(j)
			j++
		} else if i == j {
			elem = signature.In(j).Elem()
		}

		if !isSameType(elem, transformedPayload[i].Type()) {
			return nil, false
		}

	}

	return transformedPayload, true
}

//Delete callbacks and self delete himself from bus
func (l *Listeners) Delete(index *reflect.Value) {
	l.mtx.Lock()

	delete(l.callbacks, index)

	hasToRemove := len(l.callbacks) == 0 && l.remove != nil

	l.signature = nil

	l.sliceCallbacks = nil

	l.mtx.Unlock()

	if hasToRemove {
		l.remove()
	}
}

func (l *Listeners) Add(callback *reflect.Value) {
	l.mtx.Lock()

	if _, isCallback := l.callbacks[callback]; !isCallback {
		l.callbacks[callback] = callback
		l.sliceCallbacks = append(l.sliceCallbacks, callback)
	}

	if l.signature == nil {
		typeValue := callback.Type()
		l.signature = &typeValue
	}

	l.mtx.Unlock()
}

func (l *Listeners) Callbacks() []*reflect.Value {
	l.mtx.RLock()
	if l.sliceCallbacks == nil {

		callbacks := make([]*reflect.Value, len(l.callbacks))
		j := 0
		for i := range l.callbacks {
			callbacks[j] = i
			j++
		}

	}
	l.mtx.RUnlock()

	return l.sliceCallbacks
}
