package handler

import (
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"time"
)

type TaskSequenceStatus string

const (
	TaskSequenceQueued     TaskSequenceStatus = "queued"
	TaskSequenceInProgress TaskSequenceStatus = "in_progress"
	TaskSequenceFinished   TaskSequenceStatus = "finished"
	TaskSequenceErrored    TaskSequenceStatus = "errored"
)

type TaskSequenceExecutor struct {
	Name string // the source of the service that sent a .started event
}

type TaskState struct {
	Name           string
	Index          int // in case we have multiple tasks with the same name in a sequence
	TriggeredEvent keptnapimodels.KeptnContextExtendedCE
	Executors      []TaskSequenceExecutor
	FinishedEvents []keptnapimodels.KeptnContextExtendedCE
	IsFinished     bool
	Result         keptnv2.ResultType
	Status         keptnv2.StatusType
}

type TaskSequenceState struct {
	Status        TaskSequenceStatus
	Timestamp     time.Time
	InputEvent    keptnapimodels.KeptnContextExtendedCE // event that triggered the task sequence
	KeptnContext  string
	EventScope    keptnv2.EventData // project, stage, service, labels
	Shipyard      keptnv2.Shipyard  // in case the shipyard changes during the task sequence execution, keep it in the state
	CurrentTask   TaskState
	PreviousTasks []TaskState
}

type TaskSequenceQueue struct {
}
