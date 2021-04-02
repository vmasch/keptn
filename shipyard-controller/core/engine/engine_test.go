package engine_test

import (
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
			return nil, nil
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
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		ID:                 "my-id",
		Shkeptncontext:     "my-context",
		Shkeptnspecversion: "0.2.0",
		Source:             common.Stringp("my-source"),
		Specversion:        "1.0",
		Type:               common.Stringp("sh.keptn.event.dev.delivery.triggered"),
	}
}

func getDeploymentTriggeredEvent() models.KeptnContextExtendedCE {

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
				event: getDeliveryTriggeredEvent(),
			},
			expectState: state.TaskSequenceExecutionState{
				Status:     state.TaskSequenceTriggered,
				Started:    time.Now(),
				InputEvent: getDeliveryTriggeredEvent(),
				Shipyard:   keptnv2.Shipyard{},
				CurrentTask: state.TaskState{
					Name:           "deploy",
					TriggeredEvent: models.KeptnContextExtendedCE{},
					Executors:      nil,
					FinishedEvents: nil,
					IsFinished:     false,
					Result:         "",
					Status:         "",
					Started:        time.Time{},
				},
				PreviousTasks: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &engine.Engine{
				State:            tt.fields.state,
				TaskSequenceRepo: tt.fields.TaskSequenceRepo,
				ShipyardRepo:     tt.fields.ShipyardRepo,
			}
			if err := e.SetState(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("SetState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
