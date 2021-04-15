package engine

import (
	"fmt"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/eventutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/core/db"
	"github.com/keptn/keptn/shipyard-controller/core/state"
	"time"
)

type Engine struct {
	State            *state.TaskSequenceExecutionState
	TaskSequenceRepo db.ITaskSequenceExecutionStateRepo
	ShipyardRepo     IShipyardRepo
	Clock            time.Time
}

func (e *Engine) SequenceTriggered(inEvent keptnapimodels.KeptnContextExtendedCE) error {
	eventScope := &keptnv2.EventData{}
	if err := keptnv2.Decode(inEvent.Data, eventScope); err != nil {
		return err
	}

	// sync with current shipyard
	shipyard, err := e.ShipyardRepo.Sync(eventScope.GetProject())
	if err != nil {
		return err
	}

	// get task sequence for event
	taskSequence, err := e.ShipyardRepo.GetTaskSequence(*inEvent.Type)
	if err != nil {
		return err
	}

	// set new state
	e.State = state.NewTaskSequenceExecutionState(inEvent, *shipyard, *taskSequence)

	// persist state
	if err := e.TaskSequenceRepo.Store(*e.State); err != nil {
		return err
	}

	// send out neccessary event
	outEvent, err := e.deriveNextEvent()
	if err != nil {
		return err
	}

	_ = outEvent

	return nil
}

func (e *Engine) TaskStarted(event keptnapimodels.KeptnContextExtendedCE) error {
	taskName, err := ExtractTaskName(*event.Type)
	if err != nil {
		return err
	}
	if event.Source == nil {
		return fmt.Errorf("event has no source")
	}

	e.State, err = e.TaskSequenceRepo.Get(event.Shkeptncontext, event.Triggeredid, taskName)
	if err != nil {
		return err
	}

	// update list of tasks
	if _, ok := e.State.Tasks[taskName]; ok {
		e.State.Tasks[taskName] = append(e.State.Tasks[taskName], state.TaskExecutor{})
	} else {
		e.State.Tasks[taskName] = []state.TaskExecutor{{ExecutorName: *event.Source, TaskName: taskName}}
	}

	return nil
}

func (e *Engine) TaskFinished(event keptnapimodels.KeptnContextExtendedCE) error {

	taskName, err := ExtractTaskName(*event.Type)
	if err != nil {
		return err
	}

	e.State, err = e.TaskSequenceRepo.Get(event.Shkeptncontext, "", taskName)
	if err != nil {
		return err
	}

	nextEvent, err := e.deriveNextEvent()
	if err != nil {
		return err
	}

	_ = nextEvent

	return nil
}

func (e *Engine) deriveNextEvent() (*keptnapimodels.KeptnContextExtendedCE, error) {
	// merge inputEvent, previous finished events and properties of next task
	var mergedPayload interface{}
	inputDataMap := map[string]interface{}{}
	if err := keptnv2.Decode(e.State.InputEvent.Data, &inputDataMap); err != nil {
		return nil, err
	}
	mergedPayload = common.Merge(mergedPayload, inputDataMap)
	for _, task := range e.State.PreviousTasks {
		for _, finishedEvent := range task.FinishedEvents {
			mergedPayload = common.Merge(mergedPayload, finishedEvent.Data)
		}
	}

	taskProperties := map[string]interface{}{}

	taskProperties[e.State.CurrentTask.TaskName] = e.State.TaskSequence.Tasks[len(e.State.PreviousTasks)].Properties // TODO: should we store the task explicitly in ts.CurrentTask?
	mergedPayload = common.Merge(mergedPayload, taskProperties)

	event, _ := eventutils.KeptnEvent(keptnv2.GetTriggeredEventType(e.State.CurrentTask.TaskName), mergedPayload).
		WithID("NEW-ID").
		Build()

	return &event, nil
}
