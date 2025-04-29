package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"recipes-api/config"
	"recipes-api/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURI := os.Getenv("DB_URI")
	config.DB, err = config.ConnectDB(dbURI)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	router := mux.NewRouter()
	collection := config.GetCollection(config.DB)

	router.HandleFunc("/recipe", handlers.CreateRecipe(collection)).Methods("POST")
	router.HandleFunc("/recipe/{recipeId}", handlers.GetRecipe(collection)).Methods("GET")
	router.HandleFunc("/recipe/{recipeId}", handlers.EditRecipe(collection)).Methods("PUT")
	router.HandleFunc("/recipe/{recipeId}", handlers.DeleteRecipe(collection)).Methods("DELETE")
	router.HandleFunc("/recipes", handlers.GetAllRecipes(collection)).Methods("GET")

	fmt.Println("Starting the application...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
