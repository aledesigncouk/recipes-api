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
	config.ConnectDB()

	router := mux.NewRouter()
	routes.RecipeRoute(router)

	fmt.Println("Starting the application...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
