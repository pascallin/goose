package goose

import (
	"context"
	"errors"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save insert or update model
func (model *Model) Save() (result interface{}, err error) {
	primaryField := model.getPrimaryField()
	if primaryField == nil {
		return nil, errors.New("no primary field in model")
	}
	record, err := model.FindOneByID(primaryField.CurrentValue)
	if err != nil && err == mongo.ErrNoDocuments {
		return primitive.NilObjectID, err
	}
	if record != nil {
		r, err := model.FindOneByIDAndUpdate(primaryField.CurrentValue, model.CurValue)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	r, err := model.InsertOne(model.CurValue)
	if err != nil {
		return nil, err
	}
	return r, nil

}

// InsertOne insert data into collection
func (model *Model) InsertOne(v interface{}) (result interface{}, err error) {
	model.wrapCreatedAt(v)
	model.wrapUpdatedAt(v)

	insertResult, err := model.collection.InsertOne(context.TODO(), v)
	if err != nil {
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// FindOneByIDAndUpdate find one and update by id, id can be objectID or string
func (model *Model) FindOneByIDAndUpdate(id interface{}, updates interface{}) (result *mongo.SingleResult, err error) {
	model.wrapUpdatedAt(updates)

	primaryField := model.getPrimaryField()
	if reflect.ValueOf(id).Kind() == reflect.String {
		id, err = primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return nil, err
		}
	}

	after := options.After

	singleResult := model.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{primaryField.BsonName: id},
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

// DeleteOneByID delete record by id
func (model *Model) DeleteOneByID(id interface{}) (result *mongo.DeleteResult, err error) {
	primaryField := model.getPrimaryField()
	if reflect.ValueOf(id).Kind() == reflect.String {
		id, err = primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return nil, err
		}
	}
	return model.collection.DeleteOne(context.TODO(), bson.M{primaryField.BsonName: id})
}

// FindOneAndUpdate find one and update by filter
func (model *Model) FindOneAndUpdate(filter interface{}, updates interface{}) (*mongo.SingleResult, error) {
	model.wrapUpdatedAt(updates)

	after := options.After
	singleResult := model.collection.FindOneAndUpdate(
		context.TODO(),
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
	return model.collection.DeleteOne(context.TODO(), filter)
}

// BulkInsert insert batch records
func (model *Model) BulkInsert(data []interface{}, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	var operations []mongo.WriteModel
	for i := range data {
		model.wrapUpdatedAt(data[i])
		operation := mongo.NewInsertOneModel()
		operation.SetDocument(data[i])
		operations = append(operations, operation)
	}
	return model.collection.BulkWrite(context.TODO(), operations, opts...)
}

// BulkUpdate update batch records, not finished
func (model *Model) BulkUpdate(filter interface{}, data []interface{}, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	// TODO: here
	// primaryField := model.getPrimaryField()
	var operations []mongo.WriteModel
	for i := range data {
		model.wrapUpdatedAt(data[i])
		operation := mongo.NewUpdateManyModel()
		// operation.SetFilter(bson.M{primaryField.BsonName: data[i][primaryField.structName]})
		operation.SetFilter(filter)
		operation.SetUpdate(bson.M{"$set": data[i]})
		operation.SetUpsert(true)
		operations = append(operations, operation)
	}
	return model.collection.BulkWrite(context.TODO(), operations, opts...)
}

// UpdateMany update batch records
func (model *Model) UpdateMany(filter interface{}, updates interface{}) (*mongo.UpdateResult, error) {
	model.wrapUpdatedAt(updates)
	return model.collection.UpdateMany(context.TODO(), filter, bson.M{"$set": updates})
}

// DeleteMany delete batch records
func (model *Model) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	return model.collection.DeleteMany(context.TODO(), filter)
}

// SoftDeleteOne soft delete single record
func (model *Model) SoftDeleteOne(filter interface{}) (*mongo.UpdateResult, error) {
	return model.collection.UpdateOne(context.TODO(), filter, bson.M{
		"$set": bson.M{
			model.modelTime.deletedAtField.BsonName: time.Now(),
		},
	})
}

// SoftDeleteMany soft delete batch record
func (model *Model) SoftDeleteMany(filter interface{}) (*mongo.UpdateResult, error) {
	return model.collection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{
			model.modelTime.deletedAtField.BsonName: time.Now(),
		},
	})
}
