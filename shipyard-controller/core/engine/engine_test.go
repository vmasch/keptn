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

func Test_Receive_StartedEvent(t *testing.T) {

	stateRepo := &fake.ITaskSequenceExecutionStateRepoMock{}
	shipyardRepo := getFakeShipyardRepo(simpleShipyard)

	engine := engine.Engine{
		State:            state.TaskSequenceExecutionState{},
		TaskSequenceRepo: stateRepo,
		ShipyardRepo:     shipyardRepo,
	}

	engineTester := EngineTester{
		Engine:    engine,
		StateRepo: stateRepo,
	}

	sequenceTriggeredEvent := eventutils.KeptnEvent("delivery.triggered", keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}).Build()

	taskStartedEvent := eventutils.KeptnEvent("deployment.started", keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}).Build()

	// TRIGGERED EVENT
	engineTester.NewEvent(sequenceTriggeredEvent)
	//TODO: verify state
	engineTester.Execute()
	//TODO: verify state
	engineTester.Persist()
	//TODO: verify state

	// STARTED EVENT
	engineTester.NewEvent(taskStartedEvent)
	//TODO: verify state
	engineTester.Execute()
	//TODO: verify state
	engineTester.Persist()
	//TODO: verify state

}

type EngineTester struct {
	Engine    engine.Engine
	StateRepo db.ITaskSequenceExecutionStateRepo
}

func (e EngineTester) NewEvent(event models.KeptnContextExtendedCE) {
	e.Engine.SetState(event)
}

func (e EngineTester) Persist() {
	e.Engine.PersistState()
}

func (e EngineTester) Execute() {
	e.Engine.ExecuteState()
}

func (e EngineTester) VerifyState(keptnContext, triggeredID, taskName string, state state.TaskSequenceExecutionState) bool {
	executionState, _ := e.StateRepo.Get(keptnContext, triggeredID, taskName)
	_ = executionState
	return true
}

//##################

//func TestEngine_SetState(t *testing.T) {
//	type fields struct {
//		state            state.TaskSequenceExecutionState
//		TaskSequenceRepo *fake.ITaskSequenceExecutionStateRepoMock
//		ShipyardRepo     *fake2.IShipyardRepoMock
//	}
//	type args struct {
//		event models.KeptnContextExtendedCE
//	}
//	tests := []struct {
//		name        string
//		fields      fields
//		args        args
//		expectState state.TaskSequenceExecutionState
//		wantErr     bool
//	}{
//		{
//			name: "start sequence",
//			fields: fields{
//				state:            state.TaskSequenceExecutionState{},
//				TaskSequenceRepo: getFakeTaskSequenceRepo(),
//				ShipyardRepo:     getFakeShipyardRepo(simpleShipyard),
//			},
//			args: args{
//				event: getEventWithPayload("sh.keptn.event.dev.delivery.triggered", map[string]interface{}{
//					"stage": "dev",
//					"labels": map[string]interface{}{
//						"foo": "bar",
//					},
//					"deployment": map[string]interface{}{
//						"deploymentstrategy": "direct",
//					},
//					"configurationChange": map[string]interface{}{
//						"values": map[string]interface{}{
//							"foo": "bar",
//						},
//					},
//				}),
//			},
//			expectState: state.TaskSequenceExecutionState{
//				Status:    state.TaskSequenceTriggered,
//				Triggered: time.Now().Round(time.Minute),
//				InputEvent: getEventWithPayload("sh.keptn.event.dev.delivery.triggered", map[string]interface{}{
//					"stage": "dev",
//					"labels": map[string]interface{}{
//						"foo": "bar",
//					},
//					"deployment": map[string]interface{}{
//						"deploymentstrategy": "direct",
//					},
//					"configurationChange": map[string]interface{}{
//						"values": map[string]interface{}{
//							"foo": "bar",
//						},
//					},
//				}),
//				TaskSequence: getDeliveryTaskSequence(),
//				CurrentTask: state.TaskExecutor{
//					TaskName: "deployment",
//					TriggeredEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
//						"stage": "dev",
//						"labels": map[string]interface{}{
//							"foo": "bar",
//						},
//						"deployment": map[string]interface{}{
//							"deploymentstrategy": "direct",
//						},
//						"configurationChange": map[string]interface{}{
//							"values": map[string]interface{}{
//								"foo": "bar",
//							},
//						},
//					}),
//					Executors:      nil,
//					FinishedEvents: nil,
//					IsFinished:     false,
//					Result:         "",
//					Status:         "",
//					Triggered:      time.Now().Round(time.Minute),
//				},
//				PreviousTasks: nil,
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			shipyard, _ := keptnv2.DecodeShipyardYAML([]byte(simpleShipyard))
//			tt.expectState.Shipyard = *shipyard
//
//			e := &engine.Engine{
//				State:            tt.fields.state,
//				TaskSequenceRepo: tt.fields.TaskSequenceRepo,
//				ShipyardRepo:     tt.fields.ShipyardRepo,
//			}
//			if err := e.SetState(tt.args.event); (err != nil) != tt.wantErr {
//				t.Errorf("SetState() error = %v, wantErr %v", err, tt.wantErr)
//			}
//
//			// ignore generated properties
//			e.State.Triggered = tt.expectState.Triggered // ignore timestamp
//			e.State.CurrentTask.Triggered = tt.expectState.CurrentTask.Triggered
//			e.State.CurrentTask.TriggeredEvent.ID = tt.expectState.CurrentTask.TriggeredEvent.ID
//
//			assert.Equal(t, tt.expectState, e.State)
//		})
//	}
//}
//
//func TestDeriveNextTriggeredEvent(t *testing.T) {
//	shipyard, _ := keptnv2.DecodeShipyardYAML([]byte(simpleShipyard))
//
//	type args struct {
//		ts state.TaskSequenceExecutionState
//	}
//	tests := []struct {
//		name               string
//		args               args
//		wantTriggeredEvent models.KeptnContextExtendedCE
//	}{
//		{
//			name: "send initial triggered event",
//			args: args{
//				ts: state.TaskSequenceExecutionState{
//					Status:    state.TaskSequenceTriggered,
//					Triggered: time.Time{},
//					InputEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
//						"stage": "dev",
//						"labels": map[string]interface{}{
//							"foo": "bar",
//						},
//						"configurationChange": map[string]interface{}{
//							"values": "my-values",
//						},
//					}),
//					Shipyard: *shipyard,
//					TaskSequence: keptnv2.Sequence{
//						TaskName: "",
//						Tasks: []keptnv2.Task{
//							{
//								TaskName: "deployment",
//								Properties: map[string]interface{}{
//									"deploymentstrategy": "direct",
//								},
//							},
//						},
//					},
//					CurrentTask: state.TaskExecutor{
//						TaskName:       "deployment",
//						IsFinished: false,
//						Result:     "",
//						Status:     "",
//						Triggered:  time.Time{},
//					},
//					PreviousTasks: []state.TaskExecutor{},
//				},
//			},
//			wantTriggeredEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
//				"stage": "dev",
//				"labels": map[string]interface{}{
//					"foo": "bar",
//				},
//				"deployment": map[string]interface{}{
//					"deploymentstrategy": "direct",
//				},
//				"configurationChange": map[string]interface{}{
//					"values": "my-values",
//				},
//			}),
//		},
//		{
//			name: "send triggered event with properties from previous task.finished events",
//			args: args{
//				ts: state.TaskSequenceExecutionState{
//					Status:    state.TaskSequenceTriggered,
//					Triggered: time.Time{},
//					InputEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
//						"stage": "dev",
//						"labels": map[string]interface{}{
//							"foo": "bar",
//						},
//						"configurationChange": map[string]interface{}{
//							"values": "my-values",
//						},
//					}),
//					Shipyard:     *shipyard,
//					TaskSequence: getDeliveryTaskSequence(),
//					CurrentTask: state.TaskExecutor{
//						TaskName:       "test",
//						IsFinished: false,
//						Result:     "",
//						Status:     "",
//						Triggered:  time.Time{},
//					},
//					PreviousTasks: []state.TaskExecutor{
//						{
//							TaskName: "deployment",
//							FinishedEvents: []models.KeptnContextExtendedCE{
//								getEventWithPayload(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
//									"deployment": map[string]interface{}{
//										"deploymentURILocal": "my-url",
//									},
//								}),
//							},
//							Result: keptnv2.ResultPass,
//							Status: keptnv2.StatusSucceeded,
//						},
//					},
//				},
//			},
//			wantTriggeredEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), map[string]interface{}{
//				"stage": "dev",
//				"labels": map[string]interface{}{
//					"foo": "bar",
//				},
//				"deployment": map[string]interface{}{
//					"deploymentURILocal": "my-url",
//				},
//				"configurationChange": map[string]interface{}{
//					"values": "my-values",
//				},
//				"test": map[string]interface{}{
//					"teststrategy": "functional",
//				},
//			}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := engine.DeriveNextTriggeredEvent(tt.args.ts)
//
//			gotMap := map[string]interface{}{}
//
//			if err := keptnv2.Decode(got.CurrentTask.TriggeredEvent.Data, &gotMap); err != nil {
//				t.Errorf(err.Error())
//			}
//			assert.Equal(t, tt.wantTriggeredEvent.Data, gotMap)
//		})
//	}
//}

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
