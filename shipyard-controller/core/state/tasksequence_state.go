package state

import (
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"time"
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
		Status:        TaskSequenceTriggered,
		Triggered:     time.Now(),
		InputEvent:    event,
		Shipyard:      shipyard,
		TaskSequence:  sequence,
		CurrentTask:   CurrentTask{},
		PreviousTasks: []TaskResult{},
		Tasks:         map[string][]TaskExecutor{},
	}

	if len(sequence.Tasks) > 0 {
		ts.CurrentTask = CurrentTask{
			TaskName:    sequence.Tasks[0].Name,
			TriggeredID: "NEW_ID",
		}
	}

	return &ts
}

type TaskSequenceExecutionState struct {
	Status        TaskSequenceStatus
	Triggered     time.Time
	InputEvent    keptnapimodels.KeptnContextExtendedCE // event that triggered the task sequence
	Shipyard      keptnv2.Shipyard                      // in case the shipyard changes during the task sequence execution, keep it in the state
	TaskSequence  keptnv2.Sequence                      // keep the taskSequence in the state as well, because retrieving it from the shipyard, based on the incoming event will get annoying otherwise
	CurrentTask   CurrentTask
	PreviousTasks []TaskResult
	Tasks         map[string][]TaskExecutor
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
