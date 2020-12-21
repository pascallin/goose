package goose

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultSkip  int64 = 0
	defaultLimit int64 = 20
)

type FindOption struct {
	options.FindOptions
	populateChain []string
}

func (model *Model) Limit(num int64) *Model {
	model.findOpt.SetLimit(num)
	return model
}

func (model *Model) Skip(num int64) *Model {
	model.findOpt.SetSkip(num)
	return model
}

func (model *Model) clearPagination() {
	model.findOpt.SetLimit(defaultLimit)
	model.findOpt.SetLimit(defaultSkip)
}

// FindAndCountResult data struct for FindAndCount
type FindAndCountResult struct {
	Total int64
	Data  []bson.Raw
	// Data []reflect.Value
}

// FindAndCount find data and number count
func (model *Model) FindAndCount(filter bson.M) (*FindAndCountResult, error) {
	// var result []reflect.Value
	var result []bson.Raw
	options := &options.FindOptions{
		Limit: model.findOpt.Limit,
		Skip:  model.findOpt.Skip,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		model.clearPagination()
		cancel()
	}()

	cur, err := model.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		// modelType := reflect.TypeOf(model.curValue).Elem()
		// ptr := reflect.New(modelType)
		// fmt.Println(modelType, ptr, ptr.Elem().Interface())
		// err := cur.Decode(ptr.Elem().Interface())
		// if err != nil {
		// 	log.Fatal(err)
		// 	continue
		// }
		// result = append(result, ptr)
		result = append(result, cur.Current)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	total, err := model.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &FindAndCountResult{
		Total: total,
		Data:  result,
	}, nil
}

// FindOne find data by filter
func (model *Model) FindOne(filter interface{}) (*mongo.SingleResult, error) {
	singleResult := model.collection.FindOne(context.Background(), filter)
	if singleResult.Err() != nil {
		// NOTE: skip error now for ignore "no documents in result" error
		return nil, nil
	}
	return singleResult, nil
}

// FindOneByID find data by _id
func (model *Model) FindOneByID(id string) (*mongo.SingleResult, error) {
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	singleResult := model.collection.FindOne(context.Background(), bson.M{"_id": mongoID})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	return singleResult, nil
}
