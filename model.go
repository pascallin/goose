package goose

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
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
	model.initModel()
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
			case indexTag:
				_, err := model.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
					Keys: bson.D{{bsonTags.Name, 1}},
				})
				if err != nil {
					log.Fatal(err)
					continue
				}
			case defaultTag:
				if reflect.ValueOf(model.curValue).Elem().FieldByName(typeField.Name).IsZero() {
					switch valueField.Kind() {
					case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
						iVal, err := strconv.ParseInt(tagVal, 10, 64)
						if err != nil {
							log.Fatal(err)
							continue
						}
						valueField.SetInt(iVal)
					case reflect.Float32, reflect.Float64:
						fVal, err := strconv.ParseFloat(tagVal, 64)
						if err != nil {
							log.Fatal(err)
							continue
						}
						valueField.SetFloat(fVal)
					case reflect.String:
						valueField.SetString(strings.Trim(tagVal, "'"))
					case reflect.Bool:
						bVal, err := strconv.ParseBool(tagVal)
						if err != nil {
							log.Fatal(err)
							continue
						}
						valueField.SetBool(bVal)
					}
				}
			case createdAtTag:
				model.modelTime.createdAtField = &Field{
					BsonName:        bsonTags.Name,
					StructFieldName: typeField.Name,
				}
			case updatedAtTag:
				model.modelTime.updatedAtField = &Field{
					BsonName:        bsonTags.Name,
					StructFieldName: typeField.Name,
				}
			case deletedAtTag:
				model.modelTime.deletedAtField = &Field{
					BsonName:        bsonTags.Name,
					StructFieldName: typeField.Name,
				}
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

func (model *Model) initModel() {
	if model.primaryKey == "" {
		model.primaryKey = "_id"
	}
	if model.primaryKeyValue == primitive.NilObjectID {
		model.primaryKeyValue = primitive.NewObjectID()
	}
}
