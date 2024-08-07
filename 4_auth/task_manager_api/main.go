package main

import (
	"log"
	"os"
	services "task_manager_api/data"
	"task_manager_api/router"
)

func main() {
	// sets the connection string as an OS environment variable
	os.Setenv("DB_CONNECTION_STRING", DB_URL)
	// sets the secret token as an OS environment variable
	os.Setenv("JWT_SECRET_TOKEN", JWT_SECRET_TOKEN)

	// connect to DB
	err := services.ConnectDB()
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}

	log.Println("Succesfully connected to DB")

	// initiate the router and the endpoints
	router.CreateRouter(8080)
}
