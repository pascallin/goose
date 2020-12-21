package goose

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Field struct {
	Key                string
	Index              bool
	PopulateCollection string
	PopulateField      string
}

// Model Model class
type Model struct {
	collection      *mongo.Collection
	collectionName  string
	curStruct       reflect.Type
	curValue        interface{}
	fields          []Field
	findOpt         FindOption
	primaryKey      string
	primaryKeyValue primitive.ObjectID
}

func (model *Model) getCollection() (*mongo.Collection, error) {
	if DB == nil {
		return nil, errors.New("mongo not connected yet.")
	}
	return DB.Collection(model.collectionName), nil
}

// NewModel new a Model class
func NewModel(collectionName string, curValue interface{}) *Model {
	collection := getCollection(collectionName)
	return &Model{
		collection:     collection,
		collectionName: collectionName,
		curStruct:      reflect.TypeOf(curValue),
		curValue:       curValue,
	}
}

func (model *Model) GetFields() {
	for i := 0; i < model.curStruct.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := model.curStruct.Field(i)

		//Get the field tag value
		tag := field.Tag.Get(tagName)

		//Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		args := strings.Split(tag, ",")

		for _, arg := range args {
			switch arg {
			case primaryKeyTag:
				model.primaryKey = field.Name
				model.primaryKeyValue = field.
			}
		}
	}
}

func (model *Model) PrintTags() {
	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	// u := reflect.TypeOf(user)

	//Get the type and kind of our user variable
	fmt.Println("Type: ", model.curStruct.Name())
	fmt.Println("Kind: ", model.curStruct.Kind())

	for i := 0; i < model.curStruct.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := model.curStruct.Field(i)

		//Get the field tag value
		tag := field.Tag.Get(tagName)

		fmt.Printf("%d. %v(%v), tag:'%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}
