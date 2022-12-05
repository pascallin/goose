package goose

import (
	"reflect"
	"time"
)

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
