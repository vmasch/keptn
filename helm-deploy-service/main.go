package main

import (
	"context"
	"fmt"
	"log"
	"os"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/helm-deploy-service/pkg/handler"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"

	keptn "github.com/keptn/go-utils/pkg/lib"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	go switchEvent(event)
	return nil
}

func switchEvent(event cloudevents.Event) {
	serviceName := "helm-deploy-service"
	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{
		LoggingOptions: &keptn.LoggingOpts{ServiceName: &serviceName},
	})
	if err != nil {
		l := keptn.NewLogger("", event.Context.GetID(), "helm-deploy-service")
		l.Error("failed to initialize Keptn handler: " + err.Error())
		return
	}
	l := keptn.NewLogger(keptnHandler.KeptnContext, event.Context.GetID(), "helm-deploy-service")

	handlers := []handler.Handler{handler.NewDeploymentTriggeredEventHandler(keptnHandler)}

	unhandled := true
	for _, handler := range handlers {
		if handler.IsTypeHandled(event) {
			unhandled = false
			if err := handler.Handle(event); err != nil {
				l.Error(err.Error())
			}
		}
	}

	if unhandled {
		l.Error(fmt.Sprintf("Received unexpected keptn event type %s", event.Type()))
	}
}
