package main

import (
	"context"
	"fmt"
	"log"
	"task_manager_api/Delivery/router"
	domain "task_manager_api/Domain"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Connects to the mongoDB instance specified in with the connectionString
and returns a pointer to the `mongoClient`.
*/
func ConnectDB(connectionString string) (*mongo.Client, error) {
	// obtain the db connection string from the environment variables
	if connectionString == "" {
		return nil, fmt.Errorf("error: DB connection string not found. Make sure the environment variables are set correctly")
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	client, connectionErr := mongo.Connect(context.TODO(), clientOptions)
	if connectionErr != nil {
		return nil, fmt.Errorf("error: %v", connectionErr.Error())
	}

	// ping DB client to verify connection
	pingErr := client.Ping(context.TODO(), nil)
	if pingErr != nil {
		return nil, fmt.Errorf("error: %v", pingErr.Error())
	}

	return client, nil
}

/*
Adds indicies to the `users` and `tasks` collection of the databse
*/
func CreateDBIndicies(db *mongo.Database) error {
	_, err := db.Collection(domain.CollectionTasks).Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		return fmt.Errorf("error " + err.Error())
	}

	_, err = db.Collection(domain.CollectionUsers).Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		return fmt.Errorf("\n\n Error " + err.Error())
	}

	_, err = db.Collection(domain.CollectionUsers).Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		return fmt.Errorf("\n\n Error " + err.Error())
	}

	return nil
}

/*
Verifies that all the required environment variables are present in the
configured `.env` location.
*/
func CheckEnvironmentVariables() error {
	switch {
	case viper.GetString("DB_ADDRESS") == "":
		return fmt.Errorf("error while loading .env: DB_ADDRESS not found")

	case viper.GetString("DB_NAME") == "":
		return fmt.Errorf("error while loading .env: DB_NAME not found")

	case viper.GetString("SECRET_TOKEN") == "":
		return fmt.Errorf("error while loading .env: SECRET_TOKEN not found")

	case viper.GetInt("PORT") == 0:
		return fmt.Errorf("error while loading .env: PORT not found")

	case viper.GetInt("TIMEOUT") == 0:
		return fmt.Errorf("error while loading .env: PORT not found")

	case viper.GetInt("TOKEN_LIFESPAN_MINUTES") == 0:
		return fmt.Errorf("error while loading .env: PORT not found")
	}

	return nil
}

func main() {
	// load the environment variables
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// check for the environment variables
	envErr := CheckEnvironmentVariables()
	if envErr != nil {
		log.Fatal(envErr.Error())
		return
	}

	// connect to DB
	client, err := ConnectDB(viper.GetString("DB_ADDRESS"))
	db := client.Database(viper.GetString("task_API"))
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}

	// create DB indicies
	err = CreateDBIndicies(db)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}

	log.Println("Succesfully connected to DB")

	// initiate the router and the endpoints
	router.CreateRouter(viper.GetInt("PORT"), db)
}
