package engine

import (
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
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

	_, err = e.ShipyardRepo.GetTaskSequence(*event.Type)
	if err != nil {
		return err
	}

	e.State = state.NewTaskSequenceExecutionState(event, *shipyard)

	// get first task of sequence

	// get shipyard from configuration service
	// can the event trigger a sequence - look into shipyard?
	// if yes - create a new task sequence State and SetState()

	// started?

	// finished?

	panic("implement me")
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
