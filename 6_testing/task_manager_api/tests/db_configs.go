package tests

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupTaskCollection(collection *mongo.Collection) {
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		fmt.Println("\n\n Error " + err.Error())
	}
}

func SetupUserCollection(collection *mongo.Collection) {
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		fmt.Println("\n\n Error " + err.Error())
	}

	_, err = collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		fmt.Println("\n\n Error " + err.Error())
	}
}
