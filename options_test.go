package eventbus

import (
	"reflect"
	"testing"
)

func Test_defaultOrOptions(t *testing.T) {
	type args struct {
		options []*Options
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "Extract options correctly",
			args: args{
				options: []*Options{
					{Concurrency: 3, RunnerConcurrency: 1},
					{Concurrency: 5},
				},
			},
			want: &Options{
				Concurrency:       5,
				RunnerConcurrency: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultOrOptions(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultOrOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
