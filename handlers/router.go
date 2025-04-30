package handlers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(collection *mongo.Collection) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/recipe", CreateRecipe(collection)).Methods("POST")
	router.HandleFunc("/recipe/{recipeId}", GetRecipe(collection)).Methods("GET")
	router.HandleFunc("/recipe/{recipeId}", EditRecipe(collection)).Methods("PUT")
	router.HandleFunc("/recipe/{recipeId}", DeleteRecipe(collection)).Methods("DELETE")
	router.HandleFunc("/recipes", GetAllRecipes(collection)).Methods("GET")
	return router
}
