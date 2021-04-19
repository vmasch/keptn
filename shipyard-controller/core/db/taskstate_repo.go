package db

import (
	"github.com/keptn/keptn/shipyard-controller/core/state"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/task_sequence_repo.go . ITaskSequenceExecutionStateRepo
type ITaskSequenceExecutionStateRepo interface {
	Store(state state.TaskSequenceExecutionState) error
	Get(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error)
	GetSequence(sequenceName, stageName string) (*state.TaskSequenceExecutionState, error)
}

type InMemoryTaskSequenceStateRepo struct {
	store map[string]state.TaskSequenceExecutionState
}

func NewInMemoryTaskSequenceStaeRepo() *InMemoryTaskSequenceStateRepo {
	return &InMemoryTaskSequenceStateRepo{store: map[string]state.TaskSequenceExecutionState{}}
}

func (i *InMemoryTaskSequenceStateRepo) Store(state state.TaskSequenceExecutionState) error {

	key := state.CurrentSequence.SequenceName + "-" + state.CurrentStage.StageName
	i.store[key] = state
	return nil
}

func (i *InMemoryTaskSequenceStateRepo) Get(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error) {

	for _, v := range i.store {
		if v.CurrentTask.TaskName == taskName && v.CurrentTask.TriggeredID == triggeredID {
			return &v, nil
		}
	}

	return nil, nil
}
func (i *InMemoryTaskSequenceStateRepo) GetSequence(sequenceName, stageName string) (*state.TaskSequenceExecutionState, error) {
	key := sequenceName + "-" + stageName
	if value, ok := i.store[key]; ok {
		return &value, nil
	}

	return nil, nil
}
