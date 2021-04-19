package state

import (
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type TaskSequenceStatus string

const (
	TaskSequenceQueued     TaskSequenceStatus = "queued"
	TaskSequenceTriggered  TaskSequenceStatus = "triggered"
	TaskSequenceInProgress TaskSequenceStatus = "in_progress"
	TaskSequenceFinished   TaskSequenceStatus = "finished"
	TaskSequenceErrored    TaskSequenceStatus = "errored"
)

type TaskSequenceExecutor struct {
	Name string // the source of the service that sent a .started event
}

func NewTaskSequenceExecutionState(event keptnapimodels.KeptnContextExtendedCE, shipyard keptnv2.Shipyard, sequence keptnv2.Sequence) *TaskSequenceExecutionState {

	ts := TaskSequenceExecutionState{
		Status:          TaskSequenceTriggered,
		InputEvent:      event,
		Shipyard:        shipyard,
		CurrentStage:    CurrentStage{StageName: shipyard.Spec.Stages[0].Name},
		CurrentSequence: CurrentSequence{SequenceName: shipyard.Spec.Stages[0].Sequences[0].Name},
		CurrentTask:     CurrentTask{TaskName: shipyard.Spec.Stages[0].Sequences[0].Tasks[0].Name},
		PreviousTasks:   []TaskResult{},
		TaskExecutors:   map[string][]TaskExecutor{},
	}

	if len(sequence.Tasks) > 0 {
		ts.CurrentTask = CurrentTask{
			TaskName:    sequence.Tasks[0].Name,
			TriggeredID: event.ID,
		}
	}

	return &ts
}

type TaskSequenceExecutionState struct {
	Status          TaskSequenceStatus
	InputEvent      keptnapimodels.KeptnContextExtendedCE
	Shipyard        keptnv2.Shipyard
	CurrentStage    CurrentStage
	CurrentSequence CurrentSequence
	CurrentTask     CurrentTask
	PreviousTasks   []TaskResult
	TaskExecutors   map[string][]TaskExecutor
}

type CurrentStage struct {
	StageName string
}

type CurrentSequence struct {
	SequenceName string
}

type CurrentTask struct {
	TaskName    string
	TriggeredID string
}

type TaskResult struct {
	TaskName       string
	FinishedEvents []keptnapimodels.KeptnContextExtendedCE
	Result         keptnv2.ResultType
	Status         keptnv2.StatusType
}

type TaskExecutor struct {
	ExecutorName string
	TaskName     string
}

func DeriveNextState(state *TaskSequenceExecutionState) *TaskSequenceExecutionState {

	nextState := TaskSequenceExecutionState{
		Status:          state.Status,
		InputEvent:      state.InputEvent,
		Shipyard:        state.Shipyard,
		CurrentStage:    state.CurrentStage,
		CurrentSequence: state.CurrentSequence,
		CurrentTask:     CurrentTask{},
		PreviousTasks:   state.PreviousTasks,
		TaskExecutors:   state.TaskExecutors,
	}

	nextTask := GetNextTask(state)
	if nextTask != nil {
		nextState.CurrentTask = CurrentTask{TaskName: nextTask.Name, TriggeredID: state.CurrentTask.TriggeredID}
	}

	//nextStage := GetNextStage(state)
	//nextSequence := GetNextSequence(state)
	//if nextStage != nil && nextSequence != nil {
	//	nextState := TaskSequenceExecutionState{
	//		Status:          TaskSequenceTriggered,
	//		InputEvent:      state.InputEvent,
	//		Shipyard:        state.Shipyard,
	//		CurrentStage:    CurrentStage{StageName: nextStage.Name},
	//		CurrentSequence: CurrentSequence{SequenceName: nextSequence.Name},
	//		CurrentTask:     CurrentTask{},
	//		PreviousTasks:   []TaskResult{},
	//		TaskExecutors:   map[string][]TaskExecutor{},
	//	}
	//
	//	return &nextState
	//}
	return &nextState

}

func GetNextTask(state *TaskSequenceExecutionState) *keptnv2.Task {
	shipyardSpec := state.Shipyard.Spec

	for _, st := range shipyardSpec.Stages {
		if st.Name == state.CurrentStage.StageName {
			for _, seq := range st.Sequences {
				if seq.Name == state.CurrentSequence.SequenceName {
					for i, ta := range seq.Tasks {
						if ta.Name == state.CurrentTask.TaskName {
							if i+1 <= len(seq.Tasks)-1 {
								return &seq.Tasks[i+1]
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func GetNextSequence(state *TaskSequenceExecutionState) *keptnv2.Sequence {
	shipyardSpec := state.Shipyard.Spec

	for _, st := range shipyardSpec.Stages {
		if st.Name == state.CurrentStage.StageName {
			for i, seq := range st.Sequences {
				if seq.Name == state.CurrentSequence.SequenceName {
					if i+1 <= len(st.Sequences)-1 {
						return &st.Sequences[i+1]
					}
				}
			}
		}
	}
	return nil
}

func GetNextStage(state *TaskSequenceExecutionState) *keptnv2.Stage {
	shipyardSpec := state.Shipyard.Spec

	for i, st := range shipyardSpec.Stages {
		if st.Name == state.CurrentStage.StageName {
			if i+1 <= len(shipyardSpec.Stages)-1 {
				return &shipyardSpec.Stages[i+1]
			}
		}
	}
	return nil
}
