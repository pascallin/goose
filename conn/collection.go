package goose

import "go.mongodb.org/mongo-driver/mongo"

func getCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
