package main

import (
	"fmt"
	"log"
	"net/http"

	"recipes-api/config"
	"recipes-api/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to MongoDB
	config.ConnectDB()

	// Init Router
	router := mux.NewRouter()

	// Route Handlers / Endpoints
	routes.RecipeRoute(router)

	fmt.Println("Starting the application...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
