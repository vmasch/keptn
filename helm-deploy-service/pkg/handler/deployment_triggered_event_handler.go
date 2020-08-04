package handler

import (
	"fmt"

	"github.com/keptn/keptn/helm-deploy-service/pkg/helm"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"

	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnevents "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const configServiceEnv = "CONFIGURATION_SERVICE"

type DeploymentTriggeredEventHandler struct {
	keptn *keptn.Keptn
}

// NewDeploymentTriggeredEventHandler returns a new evaluation-done handler
func NewDeploymentTriggeredEventHandler(keptn *keptn.Keptn) *DeploymentTriggeredEventHandler {
	return &DeploymentTriggeredEventHandler{keptn: keptn}
}

func (DeploymentTriggeredEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.GetTriggeredEventType(keptnevents.DeploymentTaskName)
}

func (e DeploymentTriggeredEventHandler) Handle(event cloudevents.Event) error {

	triggeredEvent := &keptnevents.DeploymentTriggeredEventData{}
	if err := event.DataAs(triggeredEvent); err != nil {
		return fmt.Errorf("failed to parse EvaluationDoneEvent: %v", err)
	}

	if err := sendEvents(e.keptn, getCloudEvent(e.getDeploymentStartedEvent(triggeredEvent),
		keptnevents.GetStartedEventType(keptnevents.DeploymentTaskName), e.keptn.KeptnContext, event.ID())); err != nil {
		return err
	}

	url, err := keptn.GetServiceEndpoint(configServiceEnv)
	if err != nil {
		return err
	}

	// Read chart
	chart, err := keptnutils.GetChart(triggeredEvent.Project, triggeredEvent.Service, triggeredEvent.Stage,
		triggeredEvent.Service, url.String())
	if err != nil {
		return err
	}

	// Check for namespace
	namespace := triggeredEvent.Project + "-" + triggeredEvent.Stage
	res, err := keptnutils.ExistsNamespace(true, namespace)
	if err != nil {
		return fmt.Errorf("Failed to check if namespace %s already exists: %v", namespace, err)
	}
	if !res {
		if err := keptnutils.CreateNamespace(true, namespace); err != nil {
			return fmt.Errorf("Failed to create namespace %s: %v", namespace, err)
		}
	}

	// Execute Helm upgrade
	helmExecutor := helm.NewHelmV3Executor(e.keptn.Logger)
	if err := helmExecutor.UpgradeChart(chart, triggeredEvent.Service, namespace, nil); err != nil {
		return err
	}

	return sendEvents(e.keptn, getCloudEvent(e.getDeploymentFinishedEvent(triggeredEvent),
		keptnevents.GetFinishedEventType(keptnevents.DeploymentTaskName), e.keptn.KeptnContext, event.ID()),
		getCloudEvent(e.getDeprecatedDeploymentFinishedEvent(triggeredEvent), keptn.DeploymentFinishedEventType, e.keptn.KeptnContext, event.ID()))
}

func (e DeploymentTriggeredEventHandler) getDeploymentStartedEvent(triggeredEvent *keptnevents.DeploymentTriggeredEventData) keptnevents.DeploymentStartedEventData {
	return keptnevents.DeploymentStartedEventData{
		EventData: triggeredEvent.EventData,
	}
}

func (e DeploymentTriggeredEventHandler) getDeploymentFinishedEvent(triggeredEvent *keptnevents.DeploymentTriggeredEventData) keptnevents.DeploymentFinishedEventData {
	return keptnevents.DeploymentFinishedEventData{
		EventData:  triggeredEvent.EventData,
		Deployment: triggeredEvent.Deployment,
	}
}

func (e DeploymentTriggeredEventHandler) getDeprecatedDeploymentFinishedEvent(triggeredEvent *keptnevents.DeploymentTriggeredEventData) keptn.DeploymentFinishedEventData {
	return keptn.DeploymentFinishedEventData{
		Project:            triggeredEvent.Project,
		Stage:              triggeredEvent.Stage,
		Service:            triggeredEvent.Service,
		TestStrategy:       "performance",
		DeploymentStrategy: "blue_green_service",
		Tag:                "",
		Image:              "",
		Labels:             nil,
		DeploymentURILocal: "http://" + triggeredEvent.Service + "-canary." + triggeredEvent.Project + "-" + triggeredEvent.Stage + ".svc.cluster.local",
	}
}
