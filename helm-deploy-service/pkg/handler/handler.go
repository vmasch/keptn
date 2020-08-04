package handler

import (
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

type Handler interface {
	IsTypeHandled(event cloudevents.Event) bool
	Handle(event cloudevents.Event) error
}

func sendEvents(keptnHandler *keptnevents.Keptn, events ...cloudevents.Event) error {
	for _, outgoingEvent := range events {
		if err := keptnHandler.SendCloudEvent(outgoingEvent); err != nil {
			return err
		}
	}
	return nil
}

func getCloudEvent(data interface{}, ceType string, shkeptncontext string, triggeredID string) cloudevents.Event {

	source, _ := url.Parse("helm-deploy-service")
	contentType := "application/json"

	extensions := map[string]interface{}{"shkeptncontext": shkeptncontext}
	if triggeredID != "" {
		extensions["triggeredid"] = triggeredID
	}

	return cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        ceType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  extensions,
		}.AsV02(),
		Data: data,
	}
}
