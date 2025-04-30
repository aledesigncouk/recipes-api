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
	"recipes-api/models"
	"recipes-api/testutils"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
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
	recipe := models.Recipe{
		Name:               "Test Pizza",
		Ingredients:        []string{"Flour", "Tomato"},
		Instructions:       []string{"Mix", "Bake"},
		PrepTimeMinutes:    15,
		CookTimeMinutes:    20,
		Servings:           2,
		Difficulty:         "Easy",
		Cuisine:            "Italian",
		CaloriesPerServing: 300,
		Tags:               []string{"dinner"},
		UserID:             1,
		Image:              "pizza.jpg",
		Rating:             4.5,
		ReviewCount:        12,
		MealType:           []string{"dinner"},
	}
	body, _ := json.Marshal(recipe)

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
