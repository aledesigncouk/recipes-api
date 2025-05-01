package config

import (
	"context"
	"fmt"
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

func GetCollection(client *mongo.Client, dbName, collectionName string) (*mongo.Collection, error) {
	if client == nil {
		return nil, fmt.Errorf("mongo client is nil")
	}
	if dbName == "" || collectionName == "" {
		return nil, fmt.Errorf("dbName or collectionName is empty")
	}
	return client.Database(dbName).Collection(collectionName), nil
}

func LoadConfigFromEnv() (string, string, error) {
	dbName := os.Getenv("DB_NAME")
	collection := os.Getenv("DB_COLLECTION")

	if dbName == "" || collection == "" {
		return "", "", fmt.Errorf("DB_NAME or DB_COLLECTION env variable missing")
	}

	return dbName, collection, nil
}
