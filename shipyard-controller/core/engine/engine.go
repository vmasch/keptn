package engine

import (
	"fmt"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/core/db"
	"github.com/keptn/keptn/shipyard-controller/core/state"
	"github.com/keptn/keptn/shipyard-controller/core/utils"
	"log"
	"time"
)

type Engine struct {
	TaskSequenceRepo db.ITaskSequenceExecutionStateRepo
	ShipyardRepo     IShipyardRepo
	Clock            time.Time
}

func (e *Engine) SequenceTriggered(inEvent keptnapimodels.KeptnContextExtendedCE) error {
	eventScope := &keptnv2.EventData{}
	if err := keptnv2.Decode(inEvent.Data, eventScope); err != nil {
		return err
	}

	stage, _ := utils.ExtractStageName(*inEvent.Type)
	sequence, _ := utils.ExtractSequenceName(*inEvent.Type)
	keptnContext := inEvent.Shkeptncontext

	currentState, err := e.TaskSequenceRepo.GetBySequence(keptnContext, sequence, stage)
	if err != nil {
		return err
	}
	if currentState == nil {

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

		currentState = state.NewTaskSequenceExecutionState(inEvent, *shipyard, *taskSequence)
	}

	// persist state
	if err := e.TaskSequenceRepo.Store(*currentState); err != nil {
		return err
	}

	nextEvent, err := state.DeriveNextEvent(currentState)
	if err != nil {
		return err
	}

	log.Println("NEXT EVENT TO SEND: " + *nextEvent.Type)

	return nil
}

func (e *Engine) TaskStarted(event keptnapimodels.KeptnContextExtendedCE) error {
	taskName, err := utils.ExtractTaskName(*event.Type)
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
	if _, ok := currentState.TaskExecutors[taskName]; ok {
		currentState.TaskExecutors[taskName] = append(currentState.TaskExecutors[taskName], state.TaskExecutor{})
	} else {
		currentState.TaskExecutors[taskName] = []state.TaskExecutor{{ExecutorName: *event.Source, TaskName: taskName}}
	}

	err = e.TaskSequenceRepo.Store(*currentState)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) TaskFinished(event keptnapimodels.KeptnContextExtendedCE) error {

	finishedTaskName, err := utils.ExtractTaskName(*event.Type)
	if err != nil {
		return err
	}

	currentState, err := e.TaskSequenceRepo.Get(event.Shkeptncontext, event.Triggeredid, finishedTaskName)
	if err != nil {
		return err
	}
	_ = currentState

	// TODO: update currentState.TaskExecutors accordingly

	nextState := state.DeriveNextState(currentState)

	nextEvent, err := state.DeriveNextEvent(nextState)
	if err != nil {
		return err
	}

	e.TaskSequenceRepo.Store(*nextState)

	log.Println("NEXT EVENT TO SEND: " + *nextEvent.Type)

	return nil
}
