package util

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/non-native/mapstruct"
)

func FillStruct(data map[string]interface{}, result interface{}) {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}
}

func MapToStruct(data map[string]interface{}, result interface{}) (interface{}, error) {
	if err := mapstruct.Map2Struct(data, result); err == nil {
		return result, nil
	} else {
		return nil, err
	}
}

func InterfaceToStruct(data interface{}, result interface{}) (interface{}, error) {
	if err := mapstructure.Decode(data, result); err == nil {
		return result, err
	} else {
		return nil, err
	}
}
