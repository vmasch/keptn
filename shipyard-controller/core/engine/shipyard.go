package engine

import (
	"errors"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var ErrNoTaskSequence = errors.New("no task sequence found")

//go:generate moq -pkg fake -skip-ensure -out ./fake/shipyard.go . IShipyardRepo
type IShipyardRepo interface {
	Sync(project string) (*keptnv2.Shipyard, error)
	GetTaskSequence(eventType string) (*keptnv2.Sequence, error) // TODO: this also needs to support built-in sequences such as for evaluation
}
