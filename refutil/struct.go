package refutil

import "reflect"

func NewStructRef(value any) *StructRef {
	return &StructRef{
		value: value,
	}
}

type StructRef struct {
	value    any
	refType  reflect.Type
	refValue *reflect.Value
}

func (sr *StructRef) RefValue() *reflect.Value {
	if sr.refValue == nil {
		refValue := reflect.ValueOf(sr.value)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		sr.refValue = &refValue
	}
	return sr.refValue
}

func (sr *StructRef) RefType() reflect.Type {
	if sr.refType == nil {
		refType := reflect.TypeOf(sr.value)
		if refType.Kind() == reflect.Ptr {
			refType = refType.Elem()
		}
		sr.refType = refType
	}
	return sr.refType
}

func (sr *StructRef) Value() any {
	return sr.value
}

func (sr *StructRef) ExistField(name string) bool {
	_, ok := sr.RefType().FieldByName(name)
	return ok
}

func (sr *StructRef) GetFieldValue(name string) any {
	return sr.RefValue().FieldByName(name).Interface()
}

// SetFieldValue 设置字段值
func (sr *StructRef) SetFieldValue(name string, value any) {
	sr.RefValue().FieldByName(name).Set(reflect.ValueOf(value))
}
