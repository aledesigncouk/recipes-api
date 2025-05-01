package config_test

import (
	"context"
	"recipes-api/config"
	"recipes-api/testutils"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConnectDB(t *testing.T) {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	container, uri, _, err := testutils.StartMongoContainer(ctx, "testdb", "testcollection")
	if err != nil {
		t.Fatalf("Failed to start Mongo container: %v", err)
	}
	defer container.Terminate(ctx)

	client, err := config.ConnectDB(uri)
	assert.Equal(t, nil, err)

	collection, _ := config.GetCollection(client, "testdb", "testcollection")
	assert.Equal(t, "testdb", collection.Database().Name())
	assert.Equal(t, "testcollection", collection.Name())
	assert.NotEqual(t, nil, collection)

	// Insert test doc to verify connection works
	doc := bson.M{"title": "test recipe"}
	_, err = collection.InsertOne(ctx, doc)
	assert.Equal(t, nil, err)

	// Query back the document
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"title": "test recipe"}).Decode(&result)
	assert.Equal(t, nil, err)
	assert.Equal(t, "test recipe", result["title"])
}
