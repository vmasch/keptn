package engine_test

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/eventutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/core/db"
	"github.com/keptn/keptn/shipyard-controller/core/db/fake"
	"github.com/keptn/keptn/shipyard-controller/core/engine"
	fake2 "github.com/keptn/keptn/shipyard-controller/core/engine/fake"
	"github.com/keptn/keptn/shipyard-controller/core/state"
	"testing"
)

const simpleShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"
        - name: "delivery-direct"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"`

func getFakeShipyardRepo(shipyardContent string) *fake2.IShipyardRepoMock {
	return &fake2.IShipyardRepoMock{
		GetTaskSequenceFunc: func(eventType string) (*keptnv2.Sequence, error) {
			seq := getDeliveryTaskSequence()
			return &seq, nil
		},
		SyncFunc: func(project string) (*keptnv2.Shipyard, error) {
			return keptnv2.DecodeShipyardYAML([]byte(shipyardContent))
		},
	}
}

func getFakeTaskSequenceRepo() *fake.ITaskSequenceExecutionStateRepoMock {
	return &fake.ITaskSequenceExecutionStateRepoMock{
		GetFunc: func(keptnContext, triggeredID, taskName string) (*state.TaskSequenceExecutionState, error) {
			return nil, nil
		},
		StoreFunc: func(stateMoqParam state.TaskSequenceExecutionState) error {
			return nil
		},
		UpdateFunc: func(stateMoqParam state.TaskSequenceExecutionState) error {
			return nil
		},
	}
}

func getDeliveryTriggeredEvent() models.KeptnContextExtendedCE {
	return models.KeptnContextExtendedCE{
		Data: keptnv2.DeploymentTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "dev",
				Service: "my-service",
				Labels: map[string]string{
					"foo": "bar",
				},
			},
			ConfigurationChange: keptnv2.ConfigurationChange{Values: map[string]interface{}{
				"foo": "bar",
			}},
		},
		ID:                 "my-id",
		Shkeptncontext:     "my-context",
		Shkeptnspecversion: common.GetKeptnSpecVersion(),
		Source:             common.Stringp("my-source"),
		Specversion:        "1.0",
		Type:               common.Stringp("sh.keptn.event.dev.delivery.triggered"),
	}
}

func getEventWithPayload(eventType string, data map[string]interface{}) models.KeptnContextExtendedCE {
	data["project"] = "my-project"
	data["service"] = "my-service"
	return models.KeptnContextExtendedCE{
		Data:               data,
		ID:                 "my-id",
		Contenttype:        "application/json",
		Shkeptncontext:     "my-context",
		Shkeptnspecversion: common.GetKeptnSpecVersion(),
		Source:             common.Stringp("shipyard-controller"),
		Specversion:        "1.0",
		Type:               common.Stringp(eventType),
	}
}

func Test_ProcessTaskStartedAndFinishedEvent(t *testing.T) {

	stateRepo := &db.InMemoryTaskSequenceStateRepo{}
	shipyardRepo := getFakeShipyardRepo(simpleShipyard)

	engine := engine.Engine{
		State:            &state.TaskSequenceExecutionState{},
		TaskSequenceRepo: stateRepo,
		ShipyardRepo:     shipyardRepo,
	}

	engineTester := EngineTester{
		Engine:    engine,
		StateRepo: stateRepo,
	}

	sequenceTriggeredEvent, _ := eventutils.KeptnEvent(keptnv2.GetTriggeredEventType("deployment"), keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service"}).
		WithID("ID1").
		WithSource("cli").
		Build()

	taskStartedEvent, _ := eventutils.KeptnEvent(keptnv2.GetStartedEventType("deployment"), keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service"}).
		WithID("ID2").
		WithTriggeredID("ID1").
		WithSource("helm-service").
		Build()

	taskFinishedEvent, _ := eventutils.KeptnEvent(keptnv2.GetFinishedEventType("deployment"), keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service"}).
		WithID("ID3").
		WithTriggeredID("ID1").
		WithSource("helm-service").
		Build()

	engineTester.NewSequenceTriggeredEvent(sequenceTriggeredEvent)
	engineTester.NewTaskStartedEvent(taskStartedEvent)
	engineTester.NewTaskFinishedEvent(taskFinishedEvent)

}

type EngineTester struct {
	Engine    engine.Engine
	StateRepo db.ITaskSequenceExecutionStateRepo
}

func (e *EngineTester) SequenceTriggered(event models.KeptnContextExtendedCE) error {
	panic("implement me")
}

func (e *EngineTester) TaskStarted(event models.KeptnContextExtendedCE) error {
	panic("implement me")
}

func (e *EngineTester) TaskFinished(event models.KeptnContextExtendedCE) error {
	panic("implement me")
}

func (e *EngineTester) NewSequenceTriggeredEvent(event models.KeptnContextExtendedCE) {
	e.Engine.SequenceTriggered(event)
}

func (e *EngineTester) NewTaskStartedEvent(event models.KeptnContextExtendedCE) {
	e.Engine.TaskStarted(event)
}

func (e *EngineTester) NewTaskFinishedEvent(event models.KeptnContextExtendedCE) {
	e.Engine.TaskFinished(event)
}

func (e *EngineTester) VerifyState(keptnContext, triggeredID, taskName string, state state.TaskSequenceExecutionState) bool {
	executionState, _ := e.StateRepo.Get(keptnContext, triggeredID, taskName)
	_ = executionState
	return true
}

func getDeliveryTaskSequence() keptnv2.Sequence {
	return keptnv2.Sequence{
		Name: "delivery",
		Tasks: []keptnv2.Task{
			{
				Name: "deployment",
				Properties: map[string]interface{}{
					"deploymentstrategy": "direct",
				},
			},
			{
				Name: "test",
				Properties: map[string]interface{}{
					"teststrategy": "functional",
				},
			},
		},
	}
}
