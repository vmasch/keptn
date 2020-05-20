package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	keptnutils "github.com/keptn/go-utils/pkg/lib"

	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openTriggeredEventCollection = "open-triggered-events"

func handleTriggeredEvent(logger *keptnutils.Logger, event *models.KeptnContextExtendedCE) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(mongoDBName).Collection(openTriggeredEventCollection)

	logger.Debug("Storing event to collection " + openTriggeredEventCollection)

	eventInterface, err := transformEventToInterface(event)
	if err != nil {
		err := fmt.Errorf("failed to transform event: %v", err)
		logger.Error(err.Error())
		return err
	}

	res, err := collection.InsertOne(ctx, eventInterface)
	if err != nil {
		err := fmt.Errorf("failed to insert into collection: %v", err)
		logger.Error(err.Error())
		return err
	}
	logger.Debug(fmt.Sprintf("insertedID: %s", res.InsertedID))

	return nil
}

func deleteTriggeredEvent(eventId string) error {

	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("delete triggered event")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return err
	}

	collection := client.Database(mongoDBName).Collection(openTriggeredEventCollection)
	logger.Debug(fmt.Sprintf("Delete triggered event with ID %s from collection %s", eventId, openTriggeredEventCollection))
	if _, err := collection.DeleteOne(ctx, bson.M{"id": eventId}); err != nil {
		err := fmt.Errorf("failed to delete triggered event with ID %s from collection %s: %v", eventId, openTriggeredEventCollection, err)
		logger.Error(err.Error())
		return err
	}
	return nil
}

// GetOpenTriggeredEvents returns all open triggered events
func GetOpenTriggeredEvents(params event.GetOpenEventsParams) (*event.GetOpenEventsOKBody, error) {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("getting events from the data store")

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	searchOptions := bson.M{}
	searchOptions["type"] = params.Type

	collection := client.Database(mongoDBName).Collection(openTriggeredEventCollection)

	var result event.GetOpenEventsOKBody

	// TODO: Refactor paging functionality (duplicate code in GetEvents)
	var newNextPageKey int64
	var nextPageKey int64 = 0
	if params.NextPageKey != nil {
		tmpNextPageKey, _ := strconv.Atoi(*params.NextPageKey)
		nextPageKey = int64(tmpNextPageKey)
		newNextPageKey = nextPageKey + *params.PageSize
	} else {
		newNextPageKey = *params.PageSize
	}

	pageSize := *params.PageSize
	sortOptions := options.Find().SetSort(bson.D{{"time", -1}}).SetSkip(nextPageKey).SetLimit(pageSize)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error counting elements in events collection: %v", err))
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error finding elements in events collection: %v", err))
		return nil, err
	}
	// close the cursor after the function has completed to avoid memory leaks
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var outputEvent interface{}
		err := cur.Decode(&outputEvent)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to decode event %v", err))
			return nil, err
		}
		outputEvent, err = flattenRecursively(outputEvent, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to flatten %v", err))
			return nil, err
		}

		data, _ := json.Marshal(outputEvent)

		var keptnEvent models.KeptnContextExtendedCE
		err = keptnEvent.UnmarshalJSON(data)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to unmarshal %v", err))
			continue
		}

		result.Events = append(result.Events, &keptnEvent)
	}

	result.PageSize = pageSize
	result.TotalCount = totalCount

	if newNextPageKey < totalCount {
		result.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}

	return &result, nil
}
