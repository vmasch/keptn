package handler

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"time"
)

type IPreSendTask interface {
	Execute(t keptnv2.Task) error
}

type TaskDelay struct {
}

func (t *TaskDelay) Execute(task keptnv2.Task) error {
	if task.TriggeredAfter == "" {
		return nil
	}
	duration, err := time.ParseDuration(task.TriggeredAfter)
	if err != nil {
		return err
	}
	<-time.After(duration)
	return nil
}
