package refutil

import (
	"fmt"
	"reflect"
)

//GetMethod 获取反射方法
func GetMethod(mod interface{}, handler string) (*reflect.Value, error){
	method := reflect.ValueOf(mod).MethodByName(handler)
	if !method.IsValid(){
		return nil, fmt.Errorf("method %s is not valid", handler)
	}
	return &method, nil
}