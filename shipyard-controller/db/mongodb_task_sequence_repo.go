package db

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TaskSequenceMongoDBRepo godoc
type TaskSequenceMongoDBRepo struct {
	DbConnection MongoDBConnection
}

const taskSequenceCollectionNameSuffix = "-taskSequences"

// GetTaskSequence godoc
func (mdbrepo *TaskSequenceMongoDBRepo) GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)
	res := collection.FindOne(ctx, bson.M{"triggeredEventID": triggeredID})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Errorf("Error retrieving projects from mongoDB: %s", err.Error())
		return nil, err
	}

	taskSequenceEvent := &models.TaskSequenceEvent{}
	err = res.Decode(taskSequenceEvent)

	if err != nil {
		log.Errorf("Could not cast to *models.TaskSequenceEvent: %s", err.Error())
		return nil, err
	}

	return taskSequenceEvent, nil
}

// CreateTaskSequenceMapping godoc
func (mdbrepo *TaskSequenceMongoDBRepo) CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)

	_, err = collection.InsertOne(ctx, taskSequenceEvent)
	if err != nil {
		log.Errorf("Could not store mapping %s -> %s: %s", taskSequenceEvent.TriggeredEventID, taskSequenceEvent.TaskSequenceName, err.Error())
		return err
	}
	return nil
}

// DeleteTaskSequenceMapping godoc
func (mdbrepo *TaskSequenceMongoDBRepo) DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)

	_, err = collection.DeleteMany(ctx, bson.M{"keptnContext": keptnContext, "stage": stage, "taskSequenceName": taskSequenceName})
	if err != nil {
		log.Errorf("Could not delete entries for task %s with context %s in stage %s: %s", taskSequenceName, keptnContext, stage, err.Error())
		return err
	}
	return nil
}

// DeleteTaskSequenceCollection godoc
func (mdbrepo *TaskSequenceMongoDBRepo) DeleteTaskSequenceCollection(project string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	taskSequenceCollection := mdbrepo.getTaskSequenceCollection(project)

	if err := mdbrepo.deleteCollection(taskSequenceCollection); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (mdbrepo *TaskSequenceMongoDBRepo) deleteCollection(collection *mongo.Collection) error {
	log.Debugf("Delete collection: %s", collection.Name())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", collection.Name(), err)
		return err
	}
	return nil
}

func (mdbrepo *TaskSequenceMongoDBRepo) getTaskSequenceCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project + taskSequenceCollectionNameSuffix)
	return projectCollection
}
