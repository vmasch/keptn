package handler

import (
	"github.com/gin-gonic/gin"
	models2 "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/core/engine"
	"github.com/keptn/keptn/shipyard-controller/models"
	"log"
	"net/http"
	"strings"
)

type KeptnEventHandler struct {
	Engine engine.Engine
}

// HandleEvent godoc
// @Summary Handle event
// @Description Handle incoming cloud event
// @Tags Events
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   event     body    models.Event     true        "Event type"
// @Success 200 "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /event [post]
func (eh *KeptnEventHandler) HandleEvent(c *gin.Context) {
	event := models2.KeptnContextExtendedCE{}
	if err := c.ShouldBindJSON(event); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: common.Stringp("Invalid request format"),
		})
	}

	switch getEventStatusType(*event.Type) {
	case string(common.TriggeredEvent):
		if err := eh.Engine.SequenceTriggered(event); err != nil {
			SetInternalServerErrorResponse(err, c)
			return
		}
	case string(common.StartedEvent):
		if err := eh.Engine.TaskStarted(event); err != nil {
			SetInternalServerErrorResponse(err, c)
			return
		}
	case string(common.FinishedEvent):
		if err := eh.Engine.TaskFinished(event); err != nil {
			SetInternalServerErrorResponse(err, c)
			return
		}
	default:
		log.Println("no could not handle event")
	}

	c.Status(http.StatusOK)

}

func getEventStatusType(eventType string) string {
	split := strings.Split(eventType, ".")

	statusType := split[len(split)-1]

	return statusType
}
