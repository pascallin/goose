package goose

import (
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type Relation struct {
	from         string
	localField   string
	foreignField string
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
}

func (model *Model) getCollection() (*mongo.Collection, error) {
	if DB == nil {
		return nil, errors.New("Mongo not connected yet. ")
	}
	return DB.Collection(model.collectionName), nil
}

// NewModel new a Model class
func NewModel(collectionName string, curValue interface{}) *Model {
	collection := getCollection(collectionName)
	model := &Model{
		collection:     collection,
		collectionName: collectionName,
		curValue:       curValue,
	}
	model.getFields()
	return model
}

func (model *Model) getFields() {
	val := reflect.ValueOf(model.curValue).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get(tagName)

		//Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		for _, arg := range strings.Split(tag, ",") {
			switch arg {
			case primaryKeyTag:
				model.primaryKey = typeField.Name
				model.primaryKeyValue = valueField.Interface()
			case populateTag:
				model.relationship = append(model.relationship, Relation{
					from:       arg,
					localField: typeField.Name,
				})
			}
		}
	}
}
