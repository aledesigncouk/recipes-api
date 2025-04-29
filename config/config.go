package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB(dbURI string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	return client, nil
}

func GetCollection(client *mongo.Client) *mongo.Collection {
	collection := os.Getenv("DB_COLLECTION")
	dbName := os.Getenv("DB_NAME")

	if client == nil {
		log.Fatal("mongo client is nil")
	}

	if collection == "" || dbName == "" {
		log.Fatal("DB_COLLECTION or DB_NAME env variable missing")
	}

	return client.Database(dbName).Collection(collection)
}
