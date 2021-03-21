package eventbus

import (
	"testing"
)

func TestEvent_String(t *testing.T) {
	type fields struct {
		Name    EventName
		Payload Payload
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test format",
			fields: fields{
				Name:    "test",
				Payload: []interface{}{"test"},
			},
			want: "Event `test` \n Payload: \n[test]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Event{
				Name:    tt.fields.Name,
				Payload: tt.fields.Payload,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("Event.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
