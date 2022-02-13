package client

import (
	"context"
	"errors"
	"time"

	"github.com/praveenkumarKajla/mocketh/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
