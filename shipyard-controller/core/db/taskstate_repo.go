package db

import (
	"github.com/keptn/keptn/shipyard-controller/core/state"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/task_sequence_repo.go . ITaskSequenceExecutionStateRepo
type ITaskSequenceExecutionStateRepo interface {
	StoreTaskSequenceExecutionState(state state.TaskSequenceExecutionState) error
	GetTaskSequenceExecutionState(keptnContext, stage string) (*state.TaskSequenceExecutionState, error)
	UpdateTaskSequenceExecutionState(state state.TaskSequenceExecutionState) error
}
