package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model Model class
type Model struct {
	collection *mongo.Collection
}

// FindAndCountResult data struct for FindAndCount
type FindAndCountResult struct {
	Total int64
	Data  []bson.Raw
}

func getCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

// NewModel new a Model class
func NewModel(collectionName string) *Model {
	collection := getCollection(collectionName)
	return &Model{
		collection: collection,
	}
}

// FindAndCount find data and number count
func (Model *Model) FindAndCount(filter bson.M, pagination *Pagination) (*FindAndCountResult, error) {
	var result []bson.Raw
	pagination, err := ValidatePagination(pagination)
	if err != nil {
		return nil, err
	}
	offset := (pagination.Page - 1) * pagination.PageSize
	options := &options.FindOptions{
		Limit: &pagination.PageSize,
		Skip:  &offset,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := Model.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		result = append(result, cur.Current)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	total, err := Model.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &FindAndCountResult{
		Total: total,
		Data:  result,
	}, nil
}

// InsertOne insert data into collection
func (Model *Model) InsertOne(T interface{}) (string, error) {
	insertResult, err := Model.collection.InsertOne(context.Background(), T)
	if err != nil {
		return primitive.NilObjectID.String(), err
	}
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindOne find data by filter
func (Model *Model) FindOne(filter bson.M) (*mongo.SingleResult, error) {
	singleResult := Model.collection.
		FindOne(context.Background(), filter)
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// FindOneByID find data by _id
func (Model *Model) FindOneByID(id string) (*mongo.SingleResult, error) {
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	singleResult := Model.collection.FindOne(context.Background(), bson.M{"_id": mongoID})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// FindOneByIDAndUpdate find one and update by id
func (Model *Model) FindOneByIDAndUpdate(id string, updates bson.M) (*mongo.SingleResult, error) {
	after := options.After
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	singleResult := Model.collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": mongoID},
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
func (Model *Model) FindOneAndUpdate(filter bson.M, updates bson.M) (*mongo.SingleResult, error) {
	after := options.After
	singleResult := Model.collection.FindOneAndUpdate(
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
func (Model *Model) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	return Model.collection.DeleteOne(context.Background(), filter)
}

// DeleteOneByID delete record by id
func (Model *Model) DeleteOneByID(id string) (*mongo.DeleteResult, error) {
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return Model.collection.DeleteOne(context.Background(), bson.M{"_id": mongoID})
}

// BulkWrite insert batch records
func (Model *Model) BulkWrite(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return Model.collection.BulkWrite(context.Background(), models)
}

// UpdateMany update batch records
func (Model *Model) UpdateMany(filter bson.M, updates interface{}) (*mongo.UpdateResult, error) {
	return Model.collection.UpdateMany(context.Background(), filter, updates)
}

// DeleteMany delete batch records
func (Model *Model) DeleteMany(filter bson.M) (*mongo.DeleteResult, error) {
	return Model.collection.DeleteMany(context.Background(), filter)
}

// SoftDeleteOne soft delete single record
func (Model *Model) SoftDeleteOne(filter bson.M) (*mongo.UpdateResult, error) {
	return Model.collection.UpdateOne(context.Background(), filter, bson.M{"deletedAt": time.Now()})
}

// SoftDeleteMany soft delete batch record
func (Model *Model) SoftDeleteMany(filter bson.M) (*mongo.UpdateResult, error) {
	return Model.collection.UpdateMany(context.Background(), filter, bson.M{"deletedAt": time.Now()})
}
