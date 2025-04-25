package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"recipes-api/controllers"
	"recipes-api/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/recipe", controllers.CreateRecipe).Methods("POST")
	router.HandleFunc("/recipes", controllers.GetAllRecipes).Methods("GET")
	router.HandleFunc("/recipe/{recipeId}", controllers.GetRecipe).Methods("GET")
	router.HandleFunc("/recipe/{recipeId}", controllers.EditRecipe).Methods("PUT")
	router.HandleFunc("/recipe/{recipeId}", controllers.DeleteRecipe).Methods("DELETE")
	return router
}

func TestCreateRecipe(t *testing.T) {
	recipe := models.Recipe{
		Name:               "Test Name",
		Ingredients:        []string{"Test Ingredient 1", "Test Ingredient 2"},
		Instructions:       []string{"Test Instruction 1", "Test Instruction 2"},
		PrepTimeMinutes:    30,
		CookTimeMinutes:    45,
		Servings:           2,
		Difficulty:         "Test Difficulty",
		Cuisine:            "Test Cuisine",
		CaloriesPerServing: 200,
		Tags:               []string{"Test Tag 1", "Test Tag 2"},
		UserID:             1,
		Image:              "Test Image",
		Rating:             4.5,
		ReviewCount:        10,
		MealType:           []string{"Test Meal Type 1", "Test Meal Type 2"},
	}
	jsonValue, _ := json.Marshal(recipe)
	request, _ := http.NewRequest("POST", "/recipe", bytes.NewBuffer(jsonValue))

	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestGetAllRecipes(t *testing.T) {
	request, _ := http.NewRequest("GET", "/recipes", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUpdateRecipe(t *testing.T) {
	recipe := models.Recipe{
		Name: "Updated Name",
		Ingredients: []string{
			"Updated Ingredient 1",
			"Updated Ingredient 2",
		},
		Instructions: []string{
			"Updated Instruction 1",
			"Updated Instruction 2",
		},
		PrepTimeMinutes:    30,
		CookTimeMinutes:    45,
		Servings:           2,
		Difficulty:         "Updated Difficulty",
		Cuisine:            "Updated Cuisine",
		CaloriesPerServing: 200,
		Tags:               []string{"Updated Tag 1", "Updated Tag 2"},
		UserID:             1,
	}
	jsonValue, _ := json.Marshal(recipe)

	// You'll need to insert a test recipe first and get its ID
	// Then replace "test_id" with the actual ID
	request, _ := http.NewRequest("PUT", "/recipe/test_id", bytes.NewBuffer(jsonValue))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteRecipe(t *testing.T) {
	// You'll need to insert a test recipe first and get its ID
	// Then replace "test_id" with the actual ID
	request, _ := http.NewRequest("DELETE", "/recipe/124", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}
