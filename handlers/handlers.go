package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"recipes-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

// @Summary Create a new recipe
// @Description Adds a new recipe to the database
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body models.Recipe true "Recipe Data"
// @Success 201 {object} models.Recipe
// @Failure 400 {object} string
// @Router /recipe [post]
func CreateRecipe(collection *mongo.Collection) http.HandlerFunc {
	// TODO: consider accept and return the ID
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var recipe models.Recipe
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{"message": err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := validate.Struct(recipe); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{"message": err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		newRecipe := models.Recipe{
			// ObjectId:           primitive.NewObjectID(),
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

		// result, err := collection.InsertOne(ctx, newRecipe)
		_, err := collection.InsertOne(ctx, newRecipe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"message": err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusCreated)
		// json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(newRecipe)
	}
}

// @Summary Get a recipe
// @Description Get a recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Success 200 {object} models.Recipe
// @Failure 400 {object} string
// @Router /recipe/{recipeId} [get]
func GetRecipe(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := mux.Vars(r)
		recipeId := params["recipeId"]

		objId, err := primitive.ObjectIDFromHex(recipeId)
		if err != nil {
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
		}

		var recipe models.Recipe
		err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&recipe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(recipe)
	}
}

// @Summary Edit a recipe
// @Description Edit a recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} string
// @Router /recipe/{recipeId} [post]
func EditRecipe(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := mux.Vars(r)
		recipeId := params["recipeId"]

		objId, err := primitive.ObjectIDFromHex(recipeId)
		if err != nil {
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
		}

		var recipe models.Recipe
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if validationErr := validate.Struct(&recipe); validationErr != nil {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
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

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var updatedRecipe models.Recipe
		if result.MatchedCount == 1 {
			err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedRecipe)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedRecipe)
	}
}

// @Summary Delete a recipe
// @Description Delete a recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} string
// @Router /recipe/{recipeId} [post]
func DeleteRecipe(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := mux.Vars(r)
		recipeId := params["recipeId"]

		objId, err := primitive.ObjectIDFromHex(recipeId)
		if err != nil {
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
		}

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if result.DeletedCount < 1 {
			http.Error(w, "Recipe with specified ID not found!", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Recipe successfully deleted!"})
	}
}

// @Summary Get all recipes
// @Description Get all recipes
// @Tags recipes
// @Accept json
// @Produce json
// @Success 200 {object} models.Recipe
// @Failure 400 {object} string
// @Router /recipes [get]
func GetAllRecipes(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var recipes []models.Recipe
		for cursor.Next(ctx) {
			var recipe models.Recipe
			if err := cursor.Decode(&recipe); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			recipes = append(recipes, recipe)
		}

		if err := cursor.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(recipes)
	}
}
