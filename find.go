package goose

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindOption struct {
	options.FindOptions
	populateChain []string
}

func (model *Model) Limit() {

}
