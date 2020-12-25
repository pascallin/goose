package goose

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ref model relation for populate
type Ref struct {
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
	CurrentValue    interface{}
	IsPrimary       bool
}

// ModelTime model time
type ModelTime struct {
	createdAtField *Field
	updatedAtField *Field
	deletedAtField *Field
}

// Model Model class
type Model struct {
	collection     *mongo.Collection
	collectionName string
	CurValue       interface{}
	fields         []*Field  // only specific field will in fields
	modelTime      ModelTime // createdAt, updatedAt, deletedAt
	refs           []Ref
	findOpt        FindOption
}

func (model *Model) getCollection() *mongo.Collection {
	// if DB == nil {
	// 	return nil, errors.New("Mongo not connected yet. ")
	// }
	return DB.Collection(model.collectionName)
}

func (model *Model) getField(structName string) *Field {
	for _, field := range model.fields {
		if field.StructFieldName == structName {
			return field
		}
	}
	return nil
}

func (model *Model) getPrimaryField() *Field {
	for _, field := range model.fields {
		if field.IsPrimary {
			return field
		}
	}
	return nil
}

// NewModel new a Model class
func NewModel(collectionName string, CurValue interface{}) *Model {
	collection := getCollection(collectionName)
	model := &Model{
		collection:     collection,
		collectionName: collectionName,
		CurValue:       CurValue,
	}
	model.structTagParse()
	model.setDefault()
	return model
}

func (model *Model) setDefault() {
	primaryField := model.getPrimaryField()
	if primaryField == nil {
		model.fields = append(model.fields, &Field{
			StructFieldName: "ID",
			BsonName:        "_id",
			CurrentValue:    primitive.NewObjectID(),
			IsPrimary:       true,
		})
	}
}
