package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ORM struct {
	collection *mongo.Collection
}

type FindAndCountResult struct {
	Total int64
	Data []bson.Raw
}

func getCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

func NewORM(collectionName string) *ORM {
	collection :=  getCollection(collectionName)
	return &ORM{
		collection: collection,
	}
}

func (orm *ORM) FindAndCount (filter bson.M, pagination *Pagination) (*FindAndCountResult, error) {
	var result []bson.Raw
	pagination, err := ValidatePagination(pagination)
	if err != nil {
		return nil, err
	}
	offset := (pagination.Page - 1) * pagination.PageSize
	options := &options.FindOptions{
		Limit: &pagination.PageSize,
		Skip: &offset,
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

func (orm *ORM) InsertOne(T interface{}) (string, error) {
	insertResult, err := orm.collection.InsertOne(context.Background(), T)
	if err != nil {
		return primitive.NilObjectID.String(), err
	}
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (orm *ORM) FindOne(filter bson.M) (*mongo.SingleResult, error) {
	singleResult := orm.collection.
		FindOne(context.Background(), filter)
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

func (orm *ORM) FindOneAndUpdate(id string, filter bson.M) (*mongo.SingleResult, error) {
	after := options.After
	singleResult := orm.collection.FindOneAndUpdate(
		context.Background(),
		bson.M{ "_id": id },
		filter,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}

func (orm *ORM) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	return orm.collection.DeleteOne(context.Background(), filter)
}

func (orm *ORM) BulkWrite(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return orm.collection.BulkWrite(context.Background(), models)
}

func (orm *ORM) UpdateMany(filter bson.M, updates interface{}) (*mongo.UpdateResult, error) {
	return orm.collection.UpdateMany(context.Background(), filter, updates)
}

func (orm *ORM) DeleteMany(filter bson.M) (*mongo.DeleteResult, error){
	return orm.collection.DeleteMany(context.Background(), filter)
}

func (orm *ORM) SoftDeleteOne(filter bson.M) (*mongo.UpdateResult, error) {
	return orm.collection.UpdateOne(context.Background(), filter, bson.M{"deleted": 1})
}

func (orm *ORM) SoftDeleteMany(filter bson.M) (*mongo.UpdateResult, error) {
	return orm.collection.UpdateMany(context.Background(), filter, bson.M{"deleted": 1})
}