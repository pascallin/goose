package goose

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
)

const tagName = "goose"

const (
	// field level tags
	indexTag      = "index"
	primaryKeyTag = "primary"
	defaultTag    = "default"
	objectIDTag   = "objectID"
	// time
	createdAtTag = "createdAt"
	updatedAtTag = "updatedAt"
	deletedAtTag = "deletedAt"
	// row level tags
	refTag       = "ref"
	forignKeyTag = "forignKey"
	populateTag  = "populate"
)

func (model *Model) structTagParse() {
	val := reflect.ValueOf(model.CurValue).Elem()
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
				modelField := model.getField(typeField.Name)
				if modelField != nil {
					modelField.IsPrimary = true
				} else {
					model.fields = append(model.fields, &Field{
						BsonName:        bsonTags.Name,
						StructFieldName: typeField.Name,
						IsPrimary:       true,
					})
				}
			case indexTag:
				_, err := model.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
					Keys: bson.D{{bsonTags.Name, 1}},
				})
				if err != nil {
					log.Fatal(err)
					continue
				}
			case defaultTag:
				if reflect.ValueOf(model.CurValue).Elem().FieldByName(typeField.Name).IsZero() {
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
				model.refs = append(model.refs, Ref{
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
