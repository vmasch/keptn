package engine

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/core/db"
	"github.com/keptn/keptn/shipyard-controller/core/state"
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

	//split := strings.Split(*event.Type, ".")
	//
	//statusType := split[len(split)-1]
	//
	//eventData := &keptnv2.EventData{}
	//err := keptnv2.Decode(event.Data, eventData)
	//if err != nil {
	//	sc.logger.Error("Could not parse event data: " + err.Error())
	//	return err
	//}
	//
	//switch statusType {
	//case string(common.TriggeredEvent):
	//	return sc.handleTriggeredEvent(event)
	//case string(common.StartedEvent):
	//	return sc.handleStartedEvent(event)
	//case string(common.FinishedEvent):
	//	return sc.handleFinishedEvent(event)
	//default:
	//	return nil
	//}

	eventScope := &keptnv2.EventData{}
	if err := keptnv2.Decode(event.Data, eventScope); err != nil {
		// TODO handle error
		return err
	}
	// triggered? <stage>.<sequence>.<triggered>?
	shipyard, err := e.ShipyardRepo.Sync(eventScope.GetProject())
	if err != nil {
		return err
	}

	taskSequence, err := e.ShipyardRepo.GetTaskSequence(*event.Type)
	if err != nil {
		return err
	}

	newState, err := state.NewTaskSequenceExecutionState(event, *shipyard, *taskSequence)
	if err != nil {
		return err
	}
	e.State = *newState

	if len(e.State.TaskSequence.Tasks) > 0 {
		firstTask := e.State.TaskSequence.Tasks[0]
		e.State.CurrentTask.Name = firstTask.Name
		e.State.CurrentTask.Triggered = time.Now()
		e.State.CurrentTask.TriggeredEvent = keptnapimodels.KeptnContextExtendedCE{} // TODO: merge payload of task input event with properties of task
	}

	e.State = DeriveNextTriggeredEvent(e.State)

	// get first task of sequence

	// get shipyard from configuration service
	// can the event trigger a sequence - look into shipyard?
	// if yes - create a new task sequence State and SetState()

	// started?

	// finished?

	return nil
}

func DeriveNextTriggeredEvent(ts state.TaskSequenceExecutionState) state.TaskSequenceExecutionState {
	// merge inputEvent, previous finished events and properties of next task
	mergedPayload := ts.InputEvent.Data
	for _, task := range ts.PreviousTasks {
		for _, finishedEvent := range task.FinishedEvents {
			mergedPayload = common.Merge(mergedPayload, finishedEvent.Data)
		}
	}

	taskProperties := map[string]interface{}{}
	if len(ts.TaskSequence.Tasks) < len(ts.PreviousTasks)+1 {
		// make sure we do not have an index out of bounds access
		return ts
	}
	taskProperties[ts.CurrentTask.Name] = ts.TaskSequence.Tasks[len(ts.PreviousTasks)].Properties // TODO: should we store the task explicitly in ts.CurrentTask?
	mergedPayload = common.Merge(mergedPayload, taskProperties)

	// TODO: move this into go-utils
	source := "shipyard-controller"
	eventType := keptnv2.GetTriggeredEventType(ts.CurrentTask.Name)
	triggeredEvent := keptnapimodels.KeptnContextExtendedCE{
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               mergedPayload,
		ID:                 uuid.New().String(),
		Shkeptncontext:     ts.InputEvent.Shkeptncontext,
		Shkeptnspecversion: common.GetKeptnSpecVersion(),
		Source:             &source,
		Specversion:        cloudevents.VersionV1,
		// Time:               time.Now().String(),
		Type: &eventType,
	}

	ts.CurrentTask.TriggeredEvent = triggeredEvent
	return ts
}

func (e *Engine) ExecuteState() error {

	switch e.State.Status {
	case state.TaskSequenceTriggered:
		//e.State.CurrentTask.Name
	case state.TaskSequenceInProgress:
		break
	}

	return nil
}

func (e *Engine) PersistState() error {
	return e.TaskSequenceRepo.StoreTaskSequenceExecutionState(e.State)
}
