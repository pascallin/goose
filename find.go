package goose

import (
	"go.mongodb.org/mongo-driver/mongo/options"
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
