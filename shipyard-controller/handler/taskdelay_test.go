package handler_test

import (
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"testing"
)

func TestTaskDelay_Execute(t1 *testing.T) {
	type args struct {
		task v0_2_0.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "return error on invalid triggeredAfter string value",
			args: args{
				task: v0_2_0.Task{
					Name:           "my-task",
					TriggeredAfter: "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "wait",
			args: args{
				task: v0_2_0.Task{
					Name:           "my-task",
					TriggeredAfter: "1s",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &handler.TaskDelay{}
			if err := t.Execute(tt.args.task); (err != nil) != tt.wantErr {
				t1.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
