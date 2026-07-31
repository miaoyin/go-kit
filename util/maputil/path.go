package maputil

import (
    "errors"
    "fmt"
    "strings"
)

//GetByPathE 返回异常
func GetByPathE(m map[string]any, path string) (any, error) {
    keys := strings.Split(path, ".")
    return getNested(m, keys)
}

func GetByPath(m map[string]any, path string) any{
    val, _ := GetByPathE(m, path)
    return val
}

//GetByPathGenericE 泛型
func GetByPathGenericE[T any](m map[string]any, path string) (T, error){
    var zero T
    val, err := GetByPathE(m, path)
    if err!=nil{
        return zero, err
    }
    val2, ok := val.(T)
    if !ok {
        return zero, fmt.Errorf("invalid type value=%+v", val2)
    }
    return val2, nil
}

//GetByPathGeneric 泛型
func GetByPathGeneric[T any](m map[string]any, path string) T{
    val, _ := GetByPathGenericE[T](m, path)
    return val
}

// getNested 递归逐层查找，忽略key大小写
// 1.支持路径读a.b.c
// 2.支持嵌套map
// 3. key忽略大小写
func getNested(root map[string]any, keys []string) (any, error) {
    current := root
    for idx, rawKey := range keys {
        targetKey := strings.ToLower(rawKey)
        found := false
        var nextVal any

        // 遍历当前层所有key，小写匹配
        for k, v := range current {
            if strings.ToLower(k) == targetKey {
                found = true
                nextVal = v
                break
            }
        }
        if !found {
            return nil, fmt.Errorf("%s not found in map", targetKey)
        }

        // 最后一层，直接返回值
        if idx == len(keys)-1 {
            return nextVal, nil
        }

        // 下一层必须是 map[string]any，否则路径不合法
        subMap, ok := nextVal.(map[string]any)
        if !ok {
            return nil, errors.New("path intermediate node not map")
        }
        current = subMap
    }
    return nil, errors.New("empty path")
}
