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
	//State            *state.TaskSequenceExecutionState
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
	newState := state.NewTaskSequenceExecutionState(inEvent, *shipyard, *taskSequence)

	// persist state
	if err := e.TaskSequenceRepo.Store(*newState); err != nil {
		return err
	}

	outEvent, err := e.deriveNextEvent(newState)
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

	currentState, err := e.TaskSequenceRepo.Get(event.Shkeptncontext, event.Triggeredid, taskName)
	if err != nil {
		return err
	}

	// update list of tasks
	if _, ok := currentState.Tasks[taskName]; ok {
		currentState.Tasks[taskName] = append(currentState.Tasks[taskName], state.TaskExecutor{})
	} else {
		currentState.Tasks[taskName] = []state.TaskExecutor{{ExecutorName: *event.Source, TaskName: taskName}}
	}

	err = e.TaskSequenceRepo.Store(*currentState)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) TaskFinished(event keptnapimodels.KeptnContextExtendedCE) error {

	taskName, err := ExtractTaskName(*event.Type)
	if err != nil {
		return err
	}

	currentState, err := e.TaskSequenceRepo.Get(event.Shkeptncontext, "", taskName)
	if err != nil {
		return err
	}

	// TODO: update currentState.Tasks accordingly
	_ = currentState

	nextEvent, err := e.deriveNextEvent(currentState)
	if err != nil {
		return err
	}

	_ = nextEvent

	e.TaskSequenceRepo.Store(*currentState)

	return nil
}

func (e *Engine) deriveNextEvent(state *state.TaskSequenceExecutionState) (*keptnapimodels.KeptnContextExtendedCE, error) {
	// merge inputEvent, previous finished events and properties of next task
	var mergedPayload interface{}
	inputDataMap := map[string]interface{}{}
	if err := keptnv2.Decode(state.InputEvent.Data, &inputDataMap); err != nil {
		return nil, err
	}
	mergedPayload = common.Merge(mergedPayload, inputDataMap)
	for _, task := range state.PreviousTasks {
		for _, finishedEvent := range task.FinishedEvents {
			mergedPayload = common.Merge(mergedPayload, finishedEvent.Data)
		}
	}

	taskProperties := map[string]interface{}{}

	taskProperties[state.CurrentTask.TaskName] = state.TaskSequence.Tasks[len(state.PreviousTasks)].Properties // TODO: should we store the task explicitly in ts.CurrentTask?
	mergedPayload = common.Merge(mergedPayload, taskProperties)

	event, _ := eventutils.KeptnEvent(keptnv2.GetTriggeredEventType(state.CurrentTask.TaskName), mergedPayload).
		WithID("NEW-ID").
		Build()

	return &event, nil
}
