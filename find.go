package goose

import (
	"context"
	"log"
	"reflect"
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

// FindOption goose custom FindOption extends mongo.options.FindOption
type FindOption struct {
	options.FindOptions
	pipeline []primitive.D
}

// Limit set limit for find
func (model *Model) Limit(num int64) *Model {
	model.findOpt.SetLimit(num)
	return model
}

// Skip set skip for find
func (model *Model) Skip(num int64) *Model {
	model.findOpt.SetSkip(num)
	return model
}

// Populate populate data from other collection
func (model *Model) Populate(collectionName string) *Model {
	for _, relation := range model.refs {
		lookupStage := bson.D{
			{
				"$lookup",
				bson.D{
					{"from", relation.from},
					{"localField", relation.localField},
					{"foreignField", relation.foreignField},
					{"as", relation.as}},
			},
		}
		model.findOpt.pipeline = append(model.findOpt.pipeline, lookupStage)
	}
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
}

// FindAndCount find data and number count
func (model *Model) FindAndCount(filter bson.M) (*FindAndCountResult, error) {
	var result []bson.Raw
	options := &options.FindOptions{
		Limit: model.findOpt.Limit,
		Skip:  model.findOpt.Skip,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer func() {
		model.clearPagination()
		cancel()
	}()

	cur, err := model.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	cur.Current.Lookup()
	for cur.Next(ctx) {
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

// Find support populate find operation
func (model *Model) Find(filter interface{}) (result []bson.M, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	showLoadedCursor, err := model.collection.Aggregate(ctx, model.findOpt.pipeline)
	if err != nil {
		return nil, err
	}
	var showsLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
		return nil, err
	}
	return showsLoaded, nil
}

// FindOne find data by filter
func (model *Model) FindOne(filter interface{}) (result interface{}, err error) {
	singleResult := model.collection.FindOne(context.TODO(), filter)
	err = singleResult.Err()
	if err != nil && singleResult.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	result = singleResult.Decode(reflect.ValueOf(model.CurValue).Elem().Interface())
	return result, err
}

// FindOneByID find data by model.primaryKey
func (model *Model) FindOneByID(id interface{}) (result interface{}, err error) {
	primaryField := model.getPrimaryField()
	if reflect.ValueOf(id).Kind() == reflect.String {
		id, err = primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return nil, err
		}
	}
	singleResult := model.collection.FindOne(context.TODO(), bson.M{primaryField.BsonName: id})
	err = singleResult.Err()
	if err != nil && singleResult.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	result = singleResult.Decode(reflect.ValueOf(model.CurValue).Elem().Interface())
	return result, err
}
