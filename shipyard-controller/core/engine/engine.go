package engine

import (
	"github.com/go-openapi/strfmt"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/core/db"
	"github.com/keptn/keptn/shipyard-controller/core/state"
	"strings"
	"time"
)

type Engine struct {
	State            state.TaskSequenceExecutionState
	TaskSequenceRepo db.ITaskSequenceExecutionStateRepo
	ShipyardRepo     IShipyardRepo
	Clock            time.Time
}

type IEngine interface {
	Cancel() error
	Pause() error
	SetState(event keptnapimodels.KeptnContextExtendedCE) error
	ExecuteState() error
	PersistState() error
}

func (e *Engine) Cancel() error {
	panic("implement me")
}

func (e *Engine) Pause() error {
	panic("implement me")
}

func (e *Engine) SetState(event keptnapimodels.KeptnContextExtendedCE) error {

	statusType := getEventStatusType(*event.Type)

	switch statusType {
	case string(common.TriggeredEvent):
		return e.handleTriggeredEvent(event)
	case string(common.StartedEvent):
		return e.handleStartedEvent(event)
	case string(common.FinishedEvent):
		return e.handleFinishedEvent(event)
	default:
		return nil
	}

	return nil
}

func getEventStatusType(eventType string) string {
	split := strings.Split(eventType, ".")

	statusType := split[len(split)-1]

	return statusType
}

func DeriveTriggeredEvent(ts state.TaskSequenceExecutionState) (*keptnapimodels.KeptnContextExtendedCE, error) {
	// merge inputEvent, previous finished events and properties of next task
	var mergedPayload interface{}
	inputDataMap := map[string]interface{}{}
	if err := keptnv2.Decode(ts.InputEvent.Data, &inputDataMap); err != nil {
		return nil, err
	}
	mergedPayload = common.Merge(mergedPayload, inputDataMap)
	for _, task := range ts.PreviousTasks {
		for _, finishedEvent := range task.FinishedEvents {
			mergedPayload = common.Merge(mergedPayload, finishedEvent.Data)
		}
	}

	taskProperties := map[string]interface{}{}

	taskProperties[ts.CurrentTask.TaskName] = ts.TaskSequence.Tasks[len(ts.PreviousTasks)].Properties // TODO: should we store the task explicitly in ts.CurrentTask?
	mergedPayload = common.Merge(mergedPayload, taskProperties)

	return &keptnapimodels.KeptnContextExtendedCE{
		Contenttype:        "",
		Data:               mergedPayload,
		Extensions:         nil,
		ID:                 ts.CurrentTask.TriggeredID,
		Shkeptncontext:     "",
		Shkeptnspecversion: "",
		Source:             nil,
		Specversion:        "",
		Time:               strfmt.DateTime{},
		Triggeredid:        "",
		Type:               nil,
	}, nil

}

func (e *Engine) ExecuteState() error {
	event, err := DeriveTriggeredEvent(e.State)
	if err != nil {
		return err
	}

	// TODO: if a waitTime has been specified, wait before sending the event
	// TODO: send the event
	_ = event
	//switch e.State.Status {
	//case state.TaskSequenceTriggered:
	//	//e.State.CurrentTask.TaskName
	//case state.TaskSequenceInProgress:
	//	break
	//}

	return nil
}

func (e *Engine) PersistState() error {
	return e.TaskSequenceRepo.Store(e.State)
}

func (e *Engine) handleTriggeredEvent(event keptnapimodels.KeptnContextExtendedCE) error {
	eventScope := &keptnv2.EventData{}
	if err := keptnv2.Decode(event.Data, eventScope); err != nil {
		// TODO handle error
		return err
	}

	shipyard, err := e.ShipyardRepo.Sync(eventScope.GetProject())
	if err != nil {
		return err
	}

	taskSequence, err := e.ShipyardRepo.GetTaskSequence(*event.Type)
	if err != nil {
		return err
	}

	e.State = state.NewTaskSequenceExecutionState(event, *shipyard, *taskSequence)

	return nil
}

func (e *Engine) handleStartedEvent(event keptnapimodels.KeptnContextExtendedCE) error {

	taskName := *event.Type //TODO: this is wrong, extract task name from type

	// 1. Get current State
	// TODO: we need a lock for the sequence execution state
	currentExecutionState, err := e.TaskSequenceRepo.Get(event.Shkeptncontext, "", taskName)
	if err != nil {
		return err
	}

	// 2. Search for task and update list of tasks
	if _, ok := currentExecutionState.Tasks[taskName]; ok {
		currentExecutionState.Tasks[taskName] = append(currentExecutionState.Tasks[taskName], state.TaskExecutor{})
	} else {
		currentExecutionState.Tasks[taskName] = []state.TaskExecutor{state.TaskExecutor{}}
	}

	return nil
}

func (e *Engine) handleFinishedEvent(event keptnapimodels.KeptnContextExtendedCE) error {

	taskName := *event.Type // TODO: this is wrong, extract task name from type
	source := *event.Source

	// TODO: consider out of order events
	// 1. Get current State
	currentExecutionState, err := e.TaskSequenceRepo.Get(event.Shkeptncontext, "", taskName)
	if err != nil {
		return err
	}

	// 2. Search for task and update list of tasks
	if tasks, ok := currentExecutionState.Tasks[taskName]; ok {
		tmp := tasks[:0]
		for _, t := range tasks {
			if t.ExecutorName != source {
				tmp = append(tmp, t)
			}
		}
		currentExecutionState.Tasks[taskName] = tmp
	}

	// 3. if no more task is stored --> Update Sequence Execution State and check if next task is available
	// if there are no further tasks in the sequence, trigger next sequence(s)

	return nil
}
