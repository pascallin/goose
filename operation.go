package goose

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (model *Model) Save() error {
	key := model.primaryKey
	value := model.primaryKeyValue
	fmt.Println(key, value)
	record, err := model.FindOne(bson.M{key: value})
	if err != nil {
		return err
	}
	fmt.Println(record)
	if record != nil {
		model.FindOneAndUpdate(bson.M{key: value}, model.curValue)
	} else {
		model.InsertOne(model.curValue)
	}
	return nil
}
