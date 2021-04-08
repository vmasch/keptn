package engine_test

import (
	"github.com/bmizerany/assert"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/core/db/fake"
	"github.com/keptn/keptn/shipyard-controller/core/engine"
	fake2 "github.com/keptn/keptn/shipyard-controller/core/engine/fake"
	"github.com/keptn/keptn/shipyard-controller/core/state"
	"testing"
	"time"
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
		GetTaskSequenceExecutionStateFunc: func(keptnContext string, stage string) (*state.TaskSequenceExecutionState, error) {
			return nil, nil
		},
		StoreTaskSequenceExecutionStateFunc: func(stateMoqParam state.TaskSequenceExecutionState) error {
			return nil
		},
		UpdateTaskSequenceExecutionStateFunc: func(stateMoqParam state.TaskSequenceExecutionState) error {
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

func TestEngine_SetState(t *testing.T) {
	type fields struct {
		state            state.TaskSequenceExecutionState
		TaskSequenceRepo *fake.ITaskSequenceExecutionStateRepoMock
		ShipyardRepo     *fake2.IShipyardRepoMock
	}
	type args struct {
		event models.KeptnContextExtendedCE
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectState state.TaskSequenceExecutionState
		wantErr     bool
	}{
		{
			name: "start sequence",
			fields: fields{
				state:            state.TaskSequenceExecutionState{},
				TaskSequenceRepo: getFakeTaskSequenceRepo(),
				ShipyardRepo:     getFakeShipyardRepo(simpleShipyard),
			},
			args: args{
				event: getEventWithPayload("sh.keptn.event.dev.delivery.triggered", map[string]interface{}{
					"stage": "dev",
					"labels": map[string]interface{}{
						"foo": "bar",
					},
					"deployment": map[string]interface{}{
						"deploymentstrategy": "direct",
					},
					"configurationChange": map[string]interface{}{
						"values": map[string]interface{}{
							"foo": "bar",
						},
					},
				}),
			},
			expectState: state.TaskSequenceExecutionState{
				Status:    state.TaskSequenceTriggered,
				Triggered: time.Now().Round(time.Minute),
				InputEvent: getEventWithPayload("sh.keptn.event.dev.delivery.triggered", map[string]interface{}{
					"stage": "dev",
					"labels": map[string]interface{}{
						"foo": "bar",
					},
					"deployment": map[string]interface{}{
						"deploymentstrategy": "direct",
					},
					"configurationChange": map[string]interface{}{
						"values": map[string]interface{}{
							"foo": "bar",
						},
					},
				}),
				TaskSequence: getDeliveryTaskSequence(),
				CurrentTask: state.TaskState{
					Name: "deployment",
					TriggeredEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
						"stage": "dev",
						"labels": map[string]interface{}{
							"foo": "bar",
						},
						"deployment": map[string]interface{}{
							"deploymentstrategy": "direct",
						},
						"configurationChange": map[string]interface{}{
							"values": map[string]interface{}{
								"foo": "bar",
							},
						},
					}),
					Executors:      nil,
					FinishedEvents: nil,
					IsFinished:     false,
					Result:         "",
					Status:         "",
					Triggered:      time.Now().Round(time.Minute),
				},
				PreviousTasks: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			shipyard, _ := keptnv2.DecodeShipyardYAML([]byte(simpleShipyard))
			tt.expectState.Shipyard = *shipyard

			e := &engine.Engine{
				State:            tt.fields.state,
				TaskSequenceRepo: tt.fields.TaskSequenceRepo,
				ShipyardRepo:     tt.fields.ShipyardRepo,
			}
			if err := e.SetState(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("SetState() error = %v, wantErr %v", err, tt.wantErr)
			}

			// ignore generated properties
			e.State.Triggered = tt.expectState.Triggered // ignore timestamp
			e.State.CurrentTask.Triggered = tt.expectState.CurrentTask.Triggered
			e.State.CurrentTask.TriggeredEvent.ID = tt.expectState.CurrentTask.TriggeredEvent.ID

			assert.Equal(t, tt.expectState, e.State)
		})
	}
}

func TestDeriveNextTriggeredEvent(t *testing.T) {
	shipyard, _ := keptnv2.DecodeShipyardYAML([]byte(simpleShipyard))

	type args struct {
		ts state.TaskSequenceExecutionState
	}
	tests := []struct {
		name               string
		args               args
		wantTriggeredEvent models.KeptnContextExtendedCE
	}{
		{
			name: "send initial triggered event",
			args: args{
				ts: state.TaskSequenceExecutionState{
					Status:    state.TaskSequenceTriggered,
					Triggered: time.Time{},
					InputEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
						"stage": "dev",
						"labels": map[string]interface{}{
							"foo": "bar",
						},
						"configurationChange": map[string]interface{}{
							"values": "my-values",
						},
					}),
					Shipyard: *shipyard,
					TaskSequence: keptnv2.Sequence{
						Name: "",
						Tasks: []keptnv2.Task{
							{
								Name: "deployment",
								Properties: map[string]interface{}{
									"deploymentstrategy": "direct",
								},
							},
						},
					},
					CurrentTask: state.TaskState{
						Name:       "deployment",
						IsFinished: false,
						Result:     "",
						Status:     "",
						Triggered:  time.Time{},
					},
					PreviousTasks: []state.TaskState{},
				},
			},
			wantTriggeredEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
				"stage": "dev",
				"labels": map[string]interface{}{
					"foo": "bar",
				},
				"deployment": map[string]interface{}{
					"deploymentstrategy": "direct",
				},
				"configurationChange": map[string]interface{}{
					"values": "my-values",
				},
			}),
		},
		{
			name: "send triggered event with properties from previous task.finished events",
			args: args{
				ts: state.TaskSequenceExecutionState{
					Status:    state.TaskSequenceTriggered,
					Triggered: time.Time{},
					InputEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
						"stage": "dev",
						"labels": map[string]interface{}{
							"foo": "bar",
						},
						"configurationChange": map[string]interface{}{
							"values": "my-values",
						},
					}),
					Shipyard:     *shipyard,
					TaskSequence: getDeliveryTaskSequence(),
					CurrentTask: state.TaskState{
						Name:       "test",
						IsFinished: false,
						Result:     "",
						Status:     "",
						Triggered:  time.Time{},
					},
					PreviousTasks: []state.TaskState{
						{
							Name: "deployment",
							FinishedEvents: []models.KeptnContextExtendedCE{
								getEventWithPayload(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), map[string]interface{}{
									"deployment": map[string]interface{}{
										"deploymentURILocal": "my-url",
									},
								}),
							},
							Result: keptnv2.ResultPass,
							Status: keptnv2.StatusSucceeded,
						},
					},
				},
			},
			wantTriggeredEvent: getEventWithPayload(keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), map[string]interface{}{
				"stage": "dev",
				"labels": map[string]interface{}{
					"foo": "bar",
				},
				"deployment": map[string]interface{}{
					"deploymentURILocal": "my-url",
				},
				"configurationChange": map[string]interface{}{
					"values": "my-values",
				},
				"test": map[string]interface{}{
					"teststrategy": "functional",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.DeriveNextTriggeredEvent(tt.args.ts)

			gotMap := map[string]interface{}{}

			if err := keptnv2.Decode(got.CurrentTask.TriggeredEvent.Data, &gotMap); err != nil {
				t.Errorf(err.Error())
			}
			assert.Equal(t, tt.wantTriggeredEvent.Data, gotMap)
		})
	}
}

func getDeliveryTaskSequence() keptnv2.Sequence {
	return keptnv2.Sequence{
		Name: "",
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
