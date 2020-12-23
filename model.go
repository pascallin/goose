package goose

import (
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"

	"go.mongodb.org/mongo-driver/mongo"
)

type Relation struct {
	from         string
	localField   string
	foreignField string
	as           string
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
	model.structTagParse()
	return model
}

func (model *Model) structTagParse() {
	val := reflect.ValueOf(model.curValue).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get(tagName)

		bsonTags, err := bsoncodec.DefaultStructTagParser(typeField)
		if err != nil {
			continue
		}

		//Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		for _, arg := range strings.Split(tag, ",") {
			tagKey := arg
			var tagVal string
			if strings.Contains(string(arg), "=") {
				tagKey = strings.Split(arg, "=")[0]
				tagVal = strings.Split(arg, "=")[1]
			}
			switch tagKey {
			case primaryKeyTag:
				// model.primaryKey = typeField.Name
				model.primaryKey = bsonTags.Name
				model.primaryKeyValue = valueField.Interface()
			case populateTag:
				ref, ok := typeField.Tag.Lookup(refTag)
				if !ok {
					ref = tagVal
				}
				forignKey, ok := typeField.Tag.Lookup(forignKeyTag)
				if !ok {
					forignKey = "_id"
				}
				model.relationship = append(model.relationship, Relation{
					from: ref,
					as:   tagVal,
					// localField: typeField.Name,
					localField:   bsonTags.Name,
					foreignField: forignKey,
				})
			}
		}
	}
}
