package data

import (
	"reflect"
)

func IsNil(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}

func IsEmpty(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.ValueOf(i).Kind() {
	case reflect.Map | reflect.Array | reflect.Slice | reflect.String:
		return reflect.ValueOf(i).Len() == 0
	case reflect.Ptr | reflect.Interface:
		if reflect.ValueOf(i).IsNil() {
			return true
		}
		return IsEmpty(reflect.ValueOf(i).Elem().Interface())
	}
	return false
}
