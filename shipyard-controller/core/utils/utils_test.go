package utils

import "testing"

func TestExtractTaskName(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid event type", args{eventType: "sh.keptn.event.deployment.started"}, "deployment", false},
		{"event type with missing task", args{eventType: "sh.keptn.event.started"}, "", true},
		{"event type empty string", args{eventType: ""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTaskName(tt.args.eventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTaskName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractTaskName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
