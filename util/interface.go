package util

import (
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/non-native/mapstruct"
)

func FillStruct(data map[string]string, result map[string]string) {
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

// ValueToString returns a textual representation of the reflection value val.
// For debugging only.
func ValueToString(val reflect.Value) string {
	var str string
	if !val.IsValid() {
		return "<zero Value>"
	}
	typ := val.Type()
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		v := val
		str += typ.String()
		str += "{"
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str += ", "
			}
			str += ValueToString(v.Index(i))
		}
		str += "}"
		return str
	case reflect.Map:
		t := typ
		str = t.String()
		str += "{"
		str += "<can't iterate on maps>"
		str += "}"
		return str
	case reflect.Chan:
		str = typ.String()
		return str
	case reflect.Struct:
		t := typ
		v := val
		str += t.String()
		str += "{"
		for i, n := 0, v.NumField(); i < n; i++ {
			if i > 0 {
				str += ", "
			}
			str += ValueToString(v.Field(i))
		}
		str += "}"
		return str
	case reflect.Interface:
		return typ.String() + "(" + ValueToString(val.Elem()) + ")"
	case reflect.Func:
		v := val
		return typ.String() + "(" + strconv.FormatUint(uint64(v.Pointer()), 10) + ")"
	default:
		panic("ValueToString: can't print type " + typ.String())
	}
}
