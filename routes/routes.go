package routes

import (
	"recipes-api/controllers"

	"github.com/gorilla/mux"
)

func RecipeRoute(router *mux.Router) {
	router.HandleFunc("/recipe", controllers.CreateRecipe).Methods("POST")
	router.HandleFunc("/recipe/{recipeId}", controllers.GetRecipe).Methods("GET")
	router.HandleFunc("/recipe/{recipeId}", controllers.EditRecipe).Methods("PUT")
	router.HandleFunc("/recipe/{recipeId}", controllers.DeleteRecipe).Methods("DELETE")
	router.HandleFunc("/recipes", controllers.GetAllRecipes).Methods("GET")
}
