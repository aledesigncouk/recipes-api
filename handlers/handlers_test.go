package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"recipes-api/config"
	"recipes-api/handlers"
	"recipes-api/mockdata"
	"recipes-api/testutils"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testRouter *mux.Router

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongoC, uri, collection, err := testutils.StartMongoContainer(ctx, "testdb", "testcollection")
	if err != nil {
		log.Fatalf("Failed to start test MongoDB container: %v", err)
	}

	config.DB, _ = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	testRouter = handlers.Router(collection)

	code := m.Run()

	if err := mongoC.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate MongoDB container: %v", err)
	}

	os.Exit(code)
}

func TestCreateRecipe(t *testing.T) {
	body, _ := json.Marshal(mockdata.MockRecipe)

	req, _ := http.NewRequest("POST", "/recipe", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestGetAllRecipes(t *testing.T) {
	req, _ := http.NewRequest("GET", "/recipes", nil)
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetRecipeByID(t *testing.T) {
	inserted, err := config.DB.Database("testdb").Collection("testcollection").InsertOne(context.Background(), mockdata.MockRecipe)
	if err != nil {
		t.Fatalf("Failed to insert mock recipe: %v", err)
	}

	id := inserted.InsertedID.(primitive.ObjectID).Hex()

	req, _ := http.NewRequest("GET", "/recipe/"+id, nil)
	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}
