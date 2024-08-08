package tests

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupTaskCollection(collection *mongo.Collection) {
	collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.E{Key: "id", Value: 1}, Options: options.Index().SetUnique(true)})
}
