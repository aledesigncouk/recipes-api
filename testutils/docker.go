package testutils

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartMongoContainer(ctx context.Context, dbName string, collectionName string) (testcontainers.Container, string, *mongo.Collection, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Printf("Error starting Mongo container: %v", err)
		return nil, "", nil, err
	}

	host, _ := mongoC.Host(ctx)
	port, _ := mongoC.MappedPort(ctx, "27017")

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return nil, "", nil, err
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	mockData := []interface{}{
		bson.D{{"name", "Test Recipe 1"}, {"ingredients", []string{"Ingredient 1", "Ingredient 2"}}, {"instructions", []string{"Instruction 1", "Instruction 2"}}},
		bson.D{{"name", "Test Recipe 2"}, {"ingredients", []string{"Ingredient 3", "Ingredient 4"}}, {"instructions", []string{"Instruction 3", "Instruction 4"}}},
	}

	_, err = collection.InsertMany(ctx, mockData)
	if err != nil {
		log.Printf("Error inserting mock data: %v", err)
		return nil, "", nil, err
	}

	return mongoC, uri, collection, nil
}
