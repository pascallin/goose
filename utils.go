package goose

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func translateValue(v interface{}) (interface{}, error) {
	val := reflect.ValueOf(&v).Elem()
	fmt.Println(val, val.Kind(), val.Type())
	switch val.Kind() {
	case reflect.Int64:
		if v != 0 {
			return v, nil
		}
		return strconv.ParseInt(val.Elem().String(), 10, 64)
	case reflect.String:
		if v != "" {
			return v, nil
		}
	}
	return nil, errors.New("can not convert value")
}
