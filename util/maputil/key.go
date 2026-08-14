package maputil

import (
	"reflect"
	"strings"
)

// LowerKeys 将map中的key转为小写
func LowerKeys(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Map:
		out := reflect.MakeMap(val.Type())
		iter := val.MapRange()
		for iter.Next() {
			oldKey := iter.Key().String()
			newKey := strings.ToLower(oldKey)
			newVal := LowerKeys(iter.Value().Interface())
			out.SetMapIndex(reflect.ValueOf(newKey), reflect.ValueOf(newVal))
		}
		return out.Interface()
	default:
		return v
	}
}