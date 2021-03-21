package eventbus

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewListeners(t *testing.T) {
	if got := NewListeners(func() {}); got == nil {
		t.Errorf("NewListeners() shouldn't be nil")
	}
}

func TestListeners_Add_Delete(t *testing.T) {

	isRemoved := false

	l := NewListeners(func() {
		isRemoved = true
	})

	var (
		cb reflect.Value = reflect.ValueOf(func(payload []byte) {})
	)

	l.Add(&cb)

	if _, isVal := l.callbacks[&cb]; !isVal {
		t.Error("Listeners.Add() dosen't add reflect.Value to map")
	}

	l.Delete(&cb)

	if !isRemoved {
		t.Error("Listeners.Delete() dosen't call remove callback")
	}

	if _, isVal := l.callbacks[&cb]; isVal {
		t.Error("Listeners.Delete() dosen't remove reflect.Value from map")
	}
}

func TestListeners_Callbacks(t *testing.T) {

	l := NewListeners(func() {

	})

	var (
		cb reflect.Value = reflect.ValueOf(func(payload []byte) {})
	)
	l.Add(&cb)

	if got := l.Callbacks(); !reflect.DeepEqual(got, []*reflect.Value{&cb}) {
		t.Errorf("Listeners.Callbacks() = %v, want %v", got, []reflect.Value{cb})
	}
}

func Test_isSameType(t *testing.T) {
	type args struct {
		v1 reflect.Type
		v2 reflect.Type
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check same type",
			args: args{
				v1: reflect.ValueOf(int(1)).Type(),
				v2: reflect.ValueOf(int(1)).Type(),
			},
			want: true,
		},
		{
			name: "Check same type different",
			args: args{
				v1: reflect.ValueOf(int(1)).Type(),
				v2: reflect.ValueOf(int8(1)).Type(),
			},
		},
		{
			name: "Check same type function",
			args: args{
				v1: reflect.ValueOf(func([]int) {}).Type(),
				v2: reflect.ValueOf(func([]int) {}).Type(),
			},
			want: true,
		},
		{
			name: "Check same type empty function",
			args: args{
				v1: reflect.ValueOf(func() {}).Type(),
				v2: reflect.ValueOf(func() {}).Type(),
			},
			want: true,
		},
		{
			name: "Check same type different number of arguments for function",
			args: args{
				v1: reflect.ValueOf(func([]int) {}).Type(),
				v2: reflect.ValueOf(func() {}).Type(),
			},
			want: false,
		},
		{
			name: "Check same type different function",
			args: args{
				v1: reflect.ValueOf(func([]int) {}).Type(),
				v2: reflect.ValueOf(func([]string) {}).Type(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameType(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("isSameType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListeners_IsValidPayload(t *testing.T) {
	type fields struct {
		callbacks      map[*reflect.Value]*reflect.Value
		mtx            *sync.RWMutex
		remove         RemoveCallback
		sliceCallbacks []*reflect.Value
		signature      *reflect.Type
	}

	type args struct {
		payload []interface{}
	}

	typeV := reflect.TypeOf(func([]int) {})
	typeV2 := reflect.TypeOf(func(a string, b ...int) {})
	tests := []struct {
		name   string
		fields fields
		args   args
		want1  bool
	}{
		{
			name: "Is not valid when payload is different from signature",
			fields: fields{
				mtx:       &sync.RWMutex{},
				signature: &typeV,
			},
			args: args{
				payload: []interface{}{},
			},
			want1: false,
		},
		{
			name: "Is valid when payload is the same as signature",
			fields: fields{
				mtx:       &sync.RWMutex{},
				signature: &typeV,
			},
			args: args{
				payload: []interface{}{[]int{2}},
			},
			want1: true,
		},
		{
			name: "Is valid when payload is the same as signature with variadic parameter",
			fields: fields{
				mtx:       &sync.RWMutex{},
				signature: &typeV2,
			},
			args: args{
				payload: []interface{}{"", 2, 3},
			},
			want1: true,
		},
		{
			name: "Is valid when payload has all the manadatory parameters",
			fields: fields{
				mtx:       &sync.RWMutex{},
				signature: &typeV2,
			},
			args: args{
				payload: []interface{}{""},
			},
			want1: true,
		},
		{
			name: "Is not valid when payload dosen't include parameters which are mandatory",
			fields: fields{
				mtx:       &sync.RWMutex{},
				signature: &typeV2,
			},
			args: args{
				payload: []interface{}{},
			},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Listeners{
				callbacks:      tt.fields.callbacks,
				mtx:            tt.fields.mtx,
				remove:         tt.fields.remove,
				sliceCallbacks: tt.fields.sliceCallbacks,
				signature:      tt.fields.signature,
			}
			_, got1 := l.IsValidPayload(tt.args.payload)

			if got1 != tt.want1 {
				t.Errorf("Listeners.IsValidPayload() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestListeners_IsSame(t *testing.T) {
	type fields struct {
		callbacks      map[*reflect.Value]*reflect.Value
		mtx            *sync.RWMutex
		remove         RemoveCallback
		sliceCallbacks []*reflect.Value
		signature      *reflect.Type
	}

	type Test1 []int
	type Test2 []int

	type Test3 struct {
		t []int
	}
	type args struct {
		v *reflect.Value
	}
	t1Value := reflect.ValueOf(func(t Test1) {})
	t1 := reflect.ValueOf(func(t Test1) {}).Type()

	t2Value := reflect.ValueOf(func(t Test2) {})

	t3Value := reflect.ValueOf(func(t Test3) {})

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Same type when is no signature",
			fields: fields{
				mtx: &sync.RWMutex{},
			},
			args: args{},
			want: true,
		},

		{
			name: "Same type when is signature",
			fields: fields{
				signature: &t1,
				mtx:       &sync.RWMutex{},
			},
			args: args{
				v: &t1Value,
			},
			want: true,
		},

		{
			name: "Alias checking",
			fields: fields{
				signature: &t1,
				mtx:       &sync.RWMutex{},
			},
			args: args{
				v: &t2Value,
			},
			want: false,
		},

		{
			name: "Different type",
			fields: fields{
				signature: &t1,
				mtx:       &sync.RWMutex{},
			},
			args: args{
				v: &t3Value,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Listeners{
				callbacks:      tt.fields.callbacks,
				mtx:            tt.fields.mtx,
				remove:         tt.fields.remove,
				sliceCallbacks: tt.fields.sliceCallbacks,
				signature:      tt.fields.signature,
			}
			if got := l.IsSame(tt.args.v); got != tt.want {
				t.Errorf("Listeners.IsSame() = %v, want %v", got, tt.want)
			}
		})
	}
}
