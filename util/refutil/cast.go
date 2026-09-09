package refutil

import "reflect"

//EnsureStructPtr 确保是struct指针
func EnsureStructPtr(v any) any {
    val := reflect.ValueOf(v)

    switch val.Kind() {
    case reflect.Struct:
        //转为Struct指针
        newVal := reflect.New(val.Type())
        newVal.Elem().Set(val)
        return newVal.Interface()
    case reflect.Ptr:
        //本身是Struct指针
        if val.Elem().Kind() == reflect.Struct {
            //原始值
            return v
        }
        //多层指针
        return EnsureStructPtr(val.Elem().Interface())
    default:
        //原始值
        return v
    }
}
