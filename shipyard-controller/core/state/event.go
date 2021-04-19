package state

import (
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/eventutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
)

func DeriveNextEvent(state *TaskSequenceExecutionState) (*keptnapimodels.KeptnContextExtendedCE, error) {

	// EVENT FOR NEXT TASK
	if state.CurrentTask.TaskName != "" {
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

		//taskProperties[state.CurrentTask.TaskName] = state.CurrentSequence.Tasks[len(state.PreviousTasks)].Properties // TODO: should we store the task explicitly in ts.CurrentTask?
		mergedPayload = common.Merge(mergedPayload, taskProperties)

		event, _ := eventutils.KeptnEvent(keptnv2.GetTriggeredEventType(state.CurrentTask.TaskName), mergedPayload).
			WithID("NEW-ID").
			Build()

		return &event, nil
	}

	// EVENT FOR TRIGGERING NEXT SEQUENCE
	if state.CurrentSequence.SequenceName != "" && state.CurrentStage.StageName != "" {
		event, _ := eventutils.KeptnEvent(keptnv2.GetTriggeredEventType(state.CurrentStage.StageName+"."+state.CurrentSequence.SequenceName), keptnv2.EventData{}).
			WithID("NEW-ID").
			Build()

		return &event, nil

	}
	return nil, nil
}
