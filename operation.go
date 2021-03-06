package goose

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save insert or update model
func (model *Model) Save() error {
	key := model.primaryKey
	value := model.primaryKeyValue
	record, err := model.FindOne(bson.M{key: value})
	if err != nil {
		return err
	}
	if record != nil {
		model.FindOneAndUpdate(bson.M{key: value}, model.curValue)
	} else {
		model.InsertOne(model.curValue)
	}
	return nil
}

// InsertOne insert data into collection
func (model *Model) InsertOne(v interface{}) (string, error) {
	model.wrapCreatedAt(v)
	model.wrapUpdatedAt(v)

	data, err := bson.Marshal(v)
	if err != nil {
		return "", nil
	}

	insertResult, err := model.collection.InsertOne(context.Background(), data)
	if err != nil {
		return primitive.NilObjectID.String(), err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindOneByIDAndUpdate find one and update by id
func (model *Model) FindOneByIDAndUpdate(id string, updates interface{}) (*mongo.SingleResult, error) {
	model.wrapUpdatedAt(updates)

	after := options.After
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	singleResult := model.collection.FindOneAndUpdate(
		context.Background(),
		bson.M{model.primaryKey: mongoID},
		bson.M{
			"$set": updates,
		},
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// FindOneAndUpdate find one and update by filter
func (model *Model) FindOneAndUpdate(filter interface{}, updates interface{}) (*mongo.SingleResult, error) {
	model.wrapUpdatedAt(updates)

	after := options.After
	singleResult := model.collection.FindOneAndUpdate(
		context.Background(),
		filter,
		bson.M{
			"$set": updates,
		},
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// DeleteOne delete record by filter
func (model *Model) DeleteOne(filter interface{}) (*mongo.DeleteResult, error) {
	return model.collection.DeleteOne(context.Background(), filter)
}

// DeleteOneByID delete record by id
func (model *Model) DeleteOneByID(id string) (*mongo.DeleteResult, error) {
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return model.collection.DeleteOne(context.Background(), bson.M{model.primaryKey: mongoID})
}

// BulkWrite insert batch records
func (model *Model) BulkWrite(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	for i := range models {
		model.wrapUpdatedAt(models[i])
	}
	return model.collection.BulkWrite(context.Background(), models)
}

// UpdateMany update batch records
func (model *Model) UpdateMany(filter interface{}, updates interface{}) (*mongo.UpdateResult, error) {
	model.wrapUpdatedAt(updates)
	return model.collection.UpdateMany(context.Background(), filter, updates)
}

// DeleteMany delete batch records
func (model *Model) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	return model.collection.DeleteMany(context.Background(), filter)
}

// SoftDeleteOne soft delete single record
func (model *Model) SoftDeleteOne(filter interface{}) (*mongo.UpdateResult, error) {
	return model.collection.UpdateOne(context.Background(), filter, bson.M{
		model.modelTime.deletedAtField.BsonName: time.Now(),
	})
}

// SoftDeleteMany soft delete batch record
func (model *Model) SoftDeleteMany(filter interface{}) (*mongo.UpdateResult, error) {
	return model.collection.UpdateMany(context.Background(), filter, bson.M{
		model.modelTime.deletedAtField.BsonName: time.Now(),
	})
}

func (model *Model) wrapCreatedAt(v interface{}) {
	if model.modelTime.createdAtField != nil {
		if reflect.ValueOf(v).Elem().FieldByName(model.modelTime.createdAtField.StructFieldName).CanSet() {
			now := time.Now()
			reflect.ValueOf(v).Elem().FieldByName(model.modelTime.createdAtField.StructFieldName).Set(reflect.ValueOf(now))
		}
	}
}

func (model *Model) wrapUpdatedAt(v interface{}) {
	if model.modelTime.updatedAtField != nil {
		if reflect.ValueOf(v).Elem().FieldByName(model.modelTime.updatedAtField.StructFieldName).CanSet() {
			now := time.Now()
			reflect.ValueOf(v).Elem().FieldByName(model.modelTime.updatedAtField.StructFieldName).Set(reflect.ValueOf(now))
		}
	}
}

func (model *Model) wrapDeletedAt(v interface{}) {
	if model.modelTime.deletedAtField != nil {
		if reflect.ValueOf(v).Elem().FieldByName(model.modelTime.deletedAtField.StructFieldName).CanSet() {
			now := time.Now()
			reflect.ValueOf(v).Elem().FieldByName(model.modelTime.deletedAtField.StructFieldName).Set(reflect.ValueOf(now))
		}
	}
}
