package main

import (
	"log"
	"os"
	services "task_manager_api/data"
	"task_manager_api/router"
)

func main() {
	os.Setenv("DB_CONNECTION_STRING", DB_URL)
	err := services.ConnectDB()
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}

	log.Println("Succesfully connected to DB")
	router.CreateRouter()
}
