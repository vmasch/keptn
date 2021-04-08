package db

import (
	"github.com/keptn/keptn/shipyard-controller/core/state"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/task_sequence_repo.go . ITaskSequenceExecutionStateRepo
type ITaskSequenceExecutionStateRepo interface {
	Store(state state.TaskSequenceExecutionState) error
	Get(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error)
	Update(state state.TaskSequenceExecutionState) error
}

type InMemoryTaskSequenceStateRepo struct {
	store map[string]state.TaskSequenceExecutionState
}

func (i InMemoryTaskSequenceStateRepo) Store(state state.TaskSequenceExecutionState) error {
	panic("implement me")
}

func (i InMemoryTaskSequenceStateRepo) Get(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error) {
	panic("implement me")
}

func (i InMemoryTaskSequenceStateRepo) Update(state state.TaskSequenceExecutionState) error {
	panic("implement me")
}
