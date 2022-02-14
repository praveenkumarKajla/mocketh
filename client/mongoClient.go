package client

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/praveenkumarKajla/mocketh/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	transferCollection = "transfers"
)

var (
	onlyOnce sync.Once
)

type MongoClient struct {
	connString string
	client     *mongo.Client
}

var DBClient *MongoClient

func init() {
	mongoDBConnString := config.Config.GetString("mongoDBConnString")
	mclient, err := NewMongoClient(mongoDBConnString)
	if err != nil {
		logrus.Info("client.NewMongoClient Error", err)
		return
	}
	DBClient = mclient
}

func NewMongoClient(connString string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		connString,
	))
	if err != nil {
		logrus.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoClient{connString: connString, client: client}, nil
}

func (_mongoClient *MongoClient) GetTransferCollection() (*mongo.Collection, error) {

	collectionName := transferCollection
	collection, err := _mongoClient.GetCollection(collectionName)
	if err != nil {
		return nil, err
	}
	onlyOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mod := mongo.IndexModel{
			Keys: bson.M{
				"block_number": -1, // index in ascending order
				"block_index":  -1,
			}, Options: options.Index().SetUnique(true),
		}
		_, err = collection.Indexes().CreateOne(ctx, mod)
		if err != nil {
			logrus.Error(err)
		}
	})

	return _mongoClient.GetCollection(collectionName)
}

// Get collection handle from name
func (_mongoClient *MongoClient) GetCollection(collectionName string) (*mongo.Collection, error) {
	databaseName := config.Config.GetString("databaseName")
	if databaseName == "" {
		return nil, errors.New("empty databaseName")
	}
	if collectionName == "" {
		return nil, errors.New("empty collectionName")
	}

	collection := _mongoClient.client.Database(databaseName).Collection(collectionName)
	return collection, nil

}
