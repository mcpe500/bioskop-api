package main

import (
	"log"

	"bioskop-api/database"
	"bioskop-api/routers"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Run database migrations
	database.DBMigrate()

	r := routers.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

