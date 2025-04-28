package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	DB, err = ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
}

func ConnectDB() (*mongo.Client, error) {
	dbURI := os.Getenv("DB_URI")

	if dbURI == "" {
		return nil, fmt.Errorf("DB_URI not found")
	}

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
	collectionName := os.Getenv("DB_COLLECTION")
	dbName := os.Getenv("DB_NAME")

	// handle missing env variables

	return DB.Database(dbName).Collection(collectionName)
}
