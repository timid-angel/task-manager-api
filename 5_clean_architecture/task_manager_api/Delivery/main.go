package main

import (
	"context"
	"fmt"
	"log"
	"task_manager_api/Delivery/router"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(connectionString string) (*mongo.Client, error) {
	// obtain the db connection string from the environment variables
	if connectionString == "" {
		return nil, fmt.Errorf("Error: DB connection string not found. Make sure the environment variables are set correctly.")
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	client, connectionErr := mongo.Connect(context.TODO(), clientOptions)
	if connectionErr != nil {
		return nil, fmt.Errorf("Error: %v", connectionErr.Error())
	}

	// ping DB client to verify connection
	pingErr := client.Ping(context.TODO(), nil)
	if pingErr != nil {
		return nil, fmt.Errorf("Error: %v", pingErr.Error())
	}

	return client, nil
}

func main() {
	// load the environment variables
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	// connect to DB
	client, err := ConnectDB(viper.GetString("DB_ADDRESS"))
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}

	log.Println("Succesfully connected to DB")

	// initiate the router and the endpoints
	router.CreateRouter(viper.GetInt("PORT"), client.Database(viper.GetString("DB_NAME")))
}
