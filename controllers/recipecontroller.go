package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"recipes-api/config"
	"recipes-api/models"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var recipeCollection *mongo.Collection = config.GetCollection(config.DB)

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var recipe models.Recipe

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if validationErr := validate.Struct(&recipe); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"message": validationErr.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	newRecipe := models.Recipe{
		ObjectId:           primitive.NewObjectID(),
		ID:                 recipe.ID,
		Name:               recipe.Name,
		Ingredients:        recipe.Ingredients,
		Instructions:       recipe.Instructions,
		PrepTimeMinutes:    recipe.PrepTimeMinutes,
		CookTimeMinutes:    recipe.CookTimeMinutes,
		Servings:           recipe.Servings,
		Difficulty:         recipe.Difficulty,
		Cuisine:            recipe.Cuisine,
		CaloriesPerServing: recipe.CaloriesPerServing,
		Tags:               recipe.Tags,
		UserID:             recipe.UserID,
		Image:              recipe.Image,
		Rating:             recipe.Rating,
		ReviewCount:        recipe.ReviewCount,
		MealType:           recipe.MealType,
	}

	result, err := recipeCollection.InsertOne(ctx, newRecipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func GetRecipe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(r)
	recipeId := params["recipeId"]
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(recipeId)
	var recipe models.Recipe

	err := recipeCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&recipe)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipe)
}

func EditRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(r)
	recipeId := params["recipeId"]
	var recipe models.Recipe
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(recipeId)

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if validationErr := validate.Struct(&recipe); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"message": validationErr.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	update := bson.M{
		"name":               recipe.Name,
		"ingredients":        recipe.Ingredients,
		"instructions":       recipe.Instructions,
		"prepTimeMinutes":    recipe.PrepTimeMinutes,
		"cookTimeMinutes":    recipe.CookTimeMinutes,
		"servings":           recipe.Servings,
		"difficulty":         recipe.Difficulty,
		"cuisine":            recipe.Cuisine,
		"caloriesPerServing": recipe.CaloriesPerServing,
		"tags":               recipe.Tags,
		"userId":             recipe.UserID,
		"image":              recipe.Image,
		"rating":             recipe.Rating,
		"reviewCount":        recipe.ReviewCount,
		"mealType":           recipe.MealType,
	}

	result, err := recipeCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var updatedRecipe models.Recipe
	if result.MatchedCount == 1 {
		err := recipeCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedRecipe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"message": err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedRecipe)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	params := mux.Vars(r)
	recipeId := params["recipeId"]
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(recipeId)

	result, err := recipeCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{"message": "Recipe with specified ID not found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"message": "Recipe successfully deleted!"}
	json.NewEncoder(w).Encode(response)
}

func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var recipes []models.Recipe
	defer cancel()

	results, err := recipeCollection.Find(ctx, bson.M{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleRecipe models.Recipe
		if err = results.Decode(&singleRecipe); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"message": err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		recipes = append(recipes, singleRecipe)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipes)
}
