package goose

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

// Relation model relation for populate
type Relation struct {
	from         string
	localField   string
	foreignField string
	as           string
}

// Field model field, just using in model time for now
type Field struct {
	StructFieldName string
	BsonName        string
	DefaultValue    interface{}
}

// ModelTime model time
type ModelTime struct {
	createdAtField *Field
	updatedAtField *Field
	deletedAtField *Field
}

// Model Model class
type Model struct {
	collection      *mongo.Collection
	collectionName  string
	findOpt         FindOption
	curValue        interface{}
	primaryKey      string
	primaryKeyValue interface{}
	relationship    []Relation
	modelTime       ModelTime
}

func (model *Model) getCollection() *mongo.Collection {
	// if DB == nil {
	// 	return nil, errors.New("Mongo not connected yet. ")
	// }
	return DB.Collection(model.collectionName)
}

// NewModel new a Model class
func NewModel(collectionName string, curValue interface{}) *Model {
	collection := getCollection(collectionName)
	model := &Model{
		collection:     collection,
		collectionName: collectionName,
		curValue:       curValue,
	}
	model.structTagParse()
	model.setDefault()
	return model
}

func (model *Model) setDefault() {
	if model.primaryKey == "" {
		model.primaryKey = "_id"
	}
	if model.primaryKeyValue == primitive.NilObjectID {
		model.primaryKeyValue = primitive.NewObjectID()
	}
}
