package goose

func (model *Model) Save() {
	model.FindOneByID(model.curValue)
	model.InsertOne(model.curValue)
}
