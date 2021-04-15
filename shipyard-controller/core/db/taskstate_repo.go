package db

import (
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/core/state"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/task_sequence_repo.go . ITaskSequenceExecutionStateRepo
type ITaskSequenceExecutionStateRepo interface {
	Store(state state.TaskSequenceExecutionState) error
	Get(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error)
}

type InMemoryTaskSequenceStateRepo struct {
	store []state.TaskSequenceExecutionState
}

func (i *InMemoryTaskSequenceStateRepo) Store(state state.TaskSequenceExecutionState) error {
	i.store = append(i.store, state)
	return nil
}

func (i *InMemoryTaskSequenceStateRepo) Get(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error) {
	for _, el := range i.store {
		if el.CurrentTask.TaskName == taskName {
			return &el, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (i *InMemoryTaskSequenceStateRepo) Update(state state.TaskSequenceExecutionState) error {
	panic("implement me")
}
