package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"recipes-api/config"
	"recipes-api/handlers"

	"github.com/joho/godotenv"
)

func main() {

	var err error

	if os.Getenv("ENV") == "test" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbURI := os.Getenv("DB_URI")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	config.DB, err = config.ConnectDB(dbURI)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	dbName, collectionName, _ := config.LoadConfigFromEnv()

	collection, _ := config.GetCollection(config.DB, dbName, collectionName)
	router := handlers.Router(collection)

	fmt.Println("Starting the application...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
