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

func NewTaskSequenceExecutionState(event keptnapimodels.KeptnContextExtendedCE, shipyard keptnv2.Shipyard, sequence keptnv2.Sequence) (*TaskSequenceExecutionState, error) {
	ts := &TaskSequenceExecutionState{
		Status:        TaskSequenceTriggered,
		Triggered:     time.Now(),
		InputEvent:    event,
		Shipyard:      shipyard,
		TaskSequence:  sequence,
		CurrentTask:   TaskState{},
		PreviousTasks: []TaskState{},
	}

	return ts, nil
}

type TaskSequenceExecutionState struct {
	Status        TaskSequenceStatus
	Triggered     time.Time
	InputEvent    keptnapimodels.KeptnContextExtendedCE // event that triggered the task sequence
	Shipyard      keptnv2.Shipyard                      // in case the shipyard changes during the task sequence execution, keep it in the state
	TaskSequence  keptnv2.Sequence                      // keep the taskSequence in the state as well, because retrieving it from the shipyard, based on the incoming event will get annoying otherwise
	CurrentTask   TaskState
	PreviousTasks []TaskState // TODO: should this be a map from stage to []TaskState?
}

type TaskState struct {
	Name           string
	TriggeredEvent keptnapimodels.KeptnContextExtendedCE
	Executors      []TaskSequenceExecutor
	FinishedEvents []keptnapimodels.KeptnContextExtendedCE
	IsFinished     bool
	Result         keptnv2.ResultType
	Status         keptnv2.StatusType
	Triggered      time.Time
}

/*
type SimpleShipyardController struct {
	ConfigurationStore common.ConfigurationStore
	logger             keptncommon.LoggerInterface
}

func (sc *SimpleShipyardController) HandleIncomingEvent(event keptnapimodels.KeptnContextExtendedCE) error {
	// check if the status type is either 'triggered', 'started', or 'finished'
	split := strings.Split(*event.Type, ".")

	statusType := split[len(split)-1]

	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(event.Data, eventData)
	if err != nil {
		sc.logger.Error("Could not parse event data: " + err.Error())
		return err
	}

	switch statusType {
	case string(common.TriggeredEvent):
		return sc.handleTriggeredEvent(event)
	case string(common.StartedEvent):
		return sc.handleStartedEvent(event)
	case string(common.FinishedEvent):
		return sc.handleFinishedEvent(event)
	default:
		return nil
	}

}

func (sc *SimpleShipyardController) handleTriggeredEvent(event keptnapimodels.KeptnContextExtendedCE) error {
	if *event.Source == "shipyard-controller" {
		sc.logger.Info("Received event from myself. Ignoring.")
		return nil
	}

	// eventScope contains all properties (project, stage, service) that are needed to determine the current state within a task sequence
	// if those are not present the next action can not be determined
	eventScope, err := handler.getEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, handler.printObject(event)))

	sc.logger.Info("received event of type " + *event.Type + " from " + *event.Source)
	split := strings.Split(*event.Type, ".")

	sc.logger.Info("Checking if .triggered event should start a sequence in project " + eventScope.Project)
	// get stage and taskSequenceName - cannot tell if this is actually a task sequence triggered event though
	var stageName, taskSequenceName string
	if len(split) >= 3 {
		taskSequenceName = split[len(split)-2]
		stageName = split[len(split)-3]
	}

	resource, err := sc.ConfigurationStore.GetProjectResource(eventScope.Project, "shipyard.yaml")
	if err != nil {
		return errors.New("Could not retrieve shipyard.yaml for project " + eventScope.Project + ": " + err.Error())
	}

	shipyard, err := common.UnmarshalShipyard(resource.ResourceContent)
	if err != nil {
		// send .finished event
		return err
	}

	// validate the shipyard version - only shipyard files following the '0.2.0' spec are supported by the shipyard controller
	err = common.ValidateShipyardVersion(shipyard)
	if err != nil {
		// if the validation has not been successful: send a <task-sequence>.finished event with status=errored
		sc.logger.Error("invalid shipyard version: " + err.Error())
		// send .finished event
		return err
	}

	ts := TaskSequenceExecutionState{
		Status:        TaskSequenceTriggered,
		Triggered:       time.Now(),
		InputEvent:    event,
		KeptnContext:  "",
		Stage:         stageName,
		EventScope:    *eventScope,
		Shipyard:      *shipyard,
		CurrentTask:  TaskState{},
		PreviousTasks: []TaskState{},
	}

	if err := sc.TaskSequenceRepo.StoreTaskSequenceExecutionState(ts); err != nil {
		return err
	}
	return ts.StartExecution()
}

func getTaskSequenceInStage(stageName, taskSequenceName string, shipyard *keptnv2.Shipyard) (*keptnv2.Sequence, error) {
	for _, stage := range shipyard.Spec.Stages {
		if stage.Name == stageName {
			for _, taskSequence := range stage.Sequences {
				if taskSequence.Name == taskSequenceName {
					return &taskSequence, nil
				}
			}
			// provide built-in task sequence for evaluation
			if taskSequenceName == keptnv2.EvaluationTaskName {
				return &keptnv2.Sequence{
					Name:        "evaluation",
					TriggeredOn: nil,
					Tasks: []keptnv2.Task{
						{
							Name: keptnv2.EvaluationTaskName,
						},
					},
				}, nil
			}
			return nil, handler.errNoTaskSequence
		}
	}
	return nil, handler.errNoStage
}
*/
