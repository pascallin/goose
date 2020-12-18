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

// ORM ORM class
type ORM struct {
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

// NewORM new a orm class
func NewORM(collectionName string) *ORM {
	collection := getCollection(collectionName)
	return &ORM{
		collection: collection,
	}
}

// FindAndCount find data and number count
func (orm *ORM) FindAndCount(filter bson.M, pagination *Pagination) (*FindAndCountResult, error) {
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
	cur, err := orm.collection.Find(ctx, filter, options)
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
	total, err := orm.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &FindAndCountResult{
		Total: total,
		Data:  result,
	}, nil
}

// InsertOne insert data into collection
func (orm *ORM) InsertOne(T interface{}) (string, error) {
	insertResult, err := orm.collection.InsertOne(context.Background(), T)
	if err != nil {
		return primitive.NilObjectID.String(), err
	}
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindOne find data by filter
func (orm *ORM) FindOne(filter bson.M) (*mongo.SingleResult, error) {
	singleResult := orm.collection.
		FindOne(context.Background(), filter)
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// FindOneByIDAndUpdate find one and update by id
func (orm *ORM) FindOneByIDAndUpdate(id string, updates bson.M) (*mongo.SingleResult, error) {
	after := options.After
	singleResult := orm.collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": id},
		updates,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// FindOneAndUpdate find one and update by filter
func (orm *ORM) FindOneAndUpdate(filter bson.M, updates bson.M) (*mongo.SingleResult, error) {
	after := options.After
	singleResult := orm.collection.FindOneAndUpdate(
		context.Background(),
		filter,
		updates,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

// DeleteOne delete record by filter
func (orm *ORM) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	return orm.collection.DeleteOne(context.Background(), filter)
}

// DeleteOneByID delete record by id
func (orm *ORM) DeleteOneByID(id string) (*mongo.DeleteResult, error) {
	return orm.collection.DeleteOne(context.Background(), bson.M{"_id": id})
}

// BulkWrite insert batch records
func (orm *ORM) BulkWrite(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return orm.collection.BulkWrite(context.Background(), models)
}

// UpdateMany update batch records
func (orm *ORM) UpdateMany(filter bson.M, updates interface{}) (*mongo.UpdateResult, error) {
	return orm.collection.UpdateMany(context.Background(), filter, updates)
}

// DeleteMany delete batch records
func (orm *ORM) DeleteMany(filter bson.M) (*mongo.DeleteResult, error) {
	return orm.collection.DeleteMany(context.Background(), filter)
}

// SoftDeleteOne soft delete single record
func (orm *ORM) SoftDeleteOne(filter bson.M) (*mongo.UpdateResult, error) {
	return orm.collection.UpdateOne(context.Background(), filter, bson.M{"deletedAt": 1})
}

// SoftDeleteMany soft delete batch record
func (orm *ORM) SoftDeleteMany(filter bson.M) (*mongo.UpdateResult, error) {
	return orm.collection.UpdateMany(context.Background(), filter, bson.M{"deletedAt": 1})
}
