package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

func GetClient() *mongo.Client {
	once.Do(func() {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
	})
	return client
}