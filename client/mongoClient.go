package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/praveenkumarKajla/mocketh/config"
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
		fmt.Println("client.NewMongoClient Error", err)
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
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoClient{connString: connString, client: client}, nil
}

func (_mongoClient *MongoClient) GetCollection(databaseName string, collectionName string) (*mongo.Collection, error) {
	if databaseName == "" {
		return nil, errors.New("empty databaseName")
	}
	if collectionName == "" {
		return nil, errors.New("empty collectionName")
	}

	collection := _mongoClient.client.Database(databaseName).Collection(collectionName)
	return collection, nil

}
