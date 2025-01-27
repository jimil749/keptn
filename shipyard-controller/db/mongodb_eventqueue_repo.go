package db

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const eventQueueCollectionName = "shipyard-controller-event-queue"

// MongoDBEventQueueRepo retrieves and stores events in a mongodb collection
type MongoDBEventQueueRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBEventQueueRepo(dbConnection *MongoDBConnection) *MongoDBEventQueueRepo {
	return &MongoDBEventQueueRepo{DBConnection: dbConnection}
}

// GetQueuedEvents gets all queued events that should be sent next
func (m *MongoDBEventQueueRepo) GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error) {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DBConnection.Client.Database(getDatabaseName()).Collection(eventQueueCollectionName)

	if collection == nil {
		return nil, errors.New("invalid event type")
	}

	searchOptions := bson.M{}
	searchOptions["timestamp"] = bson.M{
		"$lte": timeutils.GetKeptnTimeStamp(timestamp),
	}

	return getQueueItemsFromCollection(collection, ctx, searchOptions)
}

func (m *MongoDBEventQueueRepo) QueueEvent(item models.QueueItem) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DBConnection.Client.Database(getDatabaseName()).Collection(eventQueueCollectionName)

	if collection == nil {
		return errors.New("invalid event type")
	}

	return insertQueueItemIntoCollection(ctx, collection, item)
}

func insertQueueItemIntoCollection(ctx context.Context, collection *mongo.Collection, item models.QueueItem) error {
	marshal, _ := json.Marshal(item)
	var eventInterface interface{}
	_ = json.Unmarshal(marshal, &eventInterface)

	existingEvent := collection.FindOne(ctx, bson.M{"eventID": item.EventID})
	if existingEvent.Err() == nil || existingEvent.Err() != mongo.ErrNoDocuments {
		return errors.New("queue item with ID " + item.EventID + " already exists in collection")
	}

	_, err := collection.InsertOne(ctx, eventInterface)
	if err != nil {
		log.Errorf("Could not insert event %s: %s", item.EventID, err.Error())
	}
	return nil
}

// DeleteQueuedEvent deletes a queue item from the collection
func (m *MongoDBEventQueueRepo) DeleteQueuedEvent(eventID string) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DBConnection.Client.Database(getDatabaseName()).Collection(eventQueueCollectionName)

	if collection == nil {
		return errors.New("invalid event type")
	}

	_, err = collection.DeleteMany(ctx, bson.M{"eventID": eventID})
	if err != nil {
		log.Errorf("Could not delete event %s : %s\n", eventID, err.Error())
		return err
	}
	log.Infof("Deleted event %s", eventID)
	return nil
}

func (m *MongoDBEventQueueRepo) IsEventInQueue(eventID string) (bool, error) {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DBConnection.Client.Database(getDatabaseName()).Collection(eventQueueCollectionName)

	if collection == nil {
		return false, errors.New("could not retrieve collection")
	}

	searchOptions := bson.M{"eventID": eventID}

	queueItems, err := getQueueItemsFromCollection(collection, ctx, searchOptions)
	if err != nil {
		return false, err
	}
	return queueItems != nil && len(queueItems) > 1, nil
}

// DeleteQueuedEvent deletes a queue item from the collection
func (m *MongoDBEventQueueRepo) DeleteQueuedEvents(scope models.EventScope) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DBConnection.Client.Database(getDatabaseName()).Collection(eventQueueCollectionName)

	if collection == nil {
		return errors.New("invalid event type")
	}

	searchOptions := bson.M{}
	if scope.KeptnContext != "" {
		searchOptions["scope.keptnContext"] = scope.KeptnContext
	}
	if scope.Project != "" {
		searchOptions["scope.project"] = scope.Project
	}
	if scope.Stage != "" {
		searchOptions["scope.stage"] = scope.Stage
	}
	if scope.Service != "" {
		searchOptions["scope.service"] = scope.Stage
	}
	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		log.Errorf("Could not delete queue items : %s\n", err.Error())
		return err
	}
	log.Info("Deleted queue items")
	return nil
}

func getQueueItemsFromCollection(collection *mongo.Collection, ctx context.Context, searchOptions bson.M, opts ...*options.FindOptions) ([]models.QueueItem, error) {
	cur, err := collection.Find(ctx, searchOptions, opts...)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, ErrNoEventFound
	} else if err != nil {
		return nil, err
	} else if cur.RemainingBatchLength() == 0 {
		return nil, ErrNoEventFound
	}

	queuedItems := []models.QueueItem{}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		queueItem := models.QueueItem{}
		err := cur.Decode(&queueItem)
		if err != nil {
			return nil, err
		}
		queuedItems = append(queuedItems, queueItem)
	}

	return queuedItems, nil
}
