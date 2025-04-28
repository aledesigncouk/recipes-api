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
	var err error
	config.DB, _ = config.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	router := mux.NewRouter()
	routes.RecipeRoute(router)

	fmt.Println("Starting the application...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
