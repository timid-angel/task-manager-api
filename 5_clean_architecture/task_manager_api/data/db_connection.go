package services

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var TaskCollection *mongo.Collection
var UserCollection *mongo.Collection
var DB_Name = "task_test" // name of the DB in the cluster

func ConnectDB() error {
	// obtain the db connection string from the environment variables
	db_string := os.Getenv("DB_CONNECTION_STRING")
	if db_string == "" {
		return ServiceError{message: "Error: DB connection string not found. Make sure the environment variables are set correctly."}
	}

	clientOptions := options.Client().ApplyURI(db_string)
	client, connectionErr := mongo.Connect(context.TODO(), clientOptions)
	if connectionErr != nil {
		return ServiceError{message: fmt.Sprintf("Error: %v", connectionErr.Error())}
	}

	// ping DB client to verify connection
	pingErr := client.Ping(context.TODO(), nil)
	if pingErr != nil {
		return ServiceError{message: fmt.Sprintf("Error: %v", pingErr.Error())}
	}

	// set the exported variables so that `task_services` can access them
	Client = client
	TaskCollection = client.Database(DB_Name).Collection("tasks")
	UserCollection = client.Database(DB_Name).Collection("users")
	return nil
}
