package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//-----------------------------------------------
//    get/set类型封装
//-----------------------------------------------

type GetCmdable interface{
	Get(ctx context.Context, key string) *redis.StringCmd
}

func GetT[T any](rdb GetCmdable, ctx context.Context, key string) (T, error){
	var zero T
	//读取
	cmd := rdb.Get(ctx, key)
	if err :=cmd.Err();err!= nil {
		return zero, err
	}
	var val any
	var err error
	switch any(zero).(type){
	case string:
		val, err = cmd.Result()
	case int:
		val, err = cmd.Int()
	case int64:
		val, err = cmd.Int64()
	case uint64:
		val, err = cmd.Uint64()
	case float32:
		val, err = cmd.Float32()
	case float64:
		val, err = cmd.Float64()
	case time.Time:
		val, err = cmd.Time()
	case bool:
		val, err = cmd.Bool()
	case []byte:
		val, err = cmd.Bytes()
	default:
		data, dataErr := cmd.Bytes()
		if dataErr!=nil{
			return zero, dataErr
		}
		err = json.Unmarshal(data, &zero)
		return zero, err
	}
	if err!=nil{
		return zero, err
	}
	return val.(T), nil
}

//GetInto 通过dest指针得到返回结果
func GetInto(rdb GetCmdable, ctx context.Context, key string, dest any) error{
	switch converted := dest.(type) {
	case *string:
		val, err := GetT[string](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *int:
		val, err := GetT[int](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *int64:
		val, err := GetT[int64](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *uint64:
		val, err := GetT[uint64](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *float32:
		val, err := GetT[float32](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *float64:
		val, err := GetT[float64](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *time.Time:
		val, err := GetT[time.Time](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *bool:
		val, err := GetT[bool](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	case *[]byte:
		val, err := GetT[[]byte](rdb, ctx, key)
		if err!=nil{
			return err
		}
		*converted = val
		return nil
	default:
		data, err := GetT[[]byte](rdb, ctx, key)
		if err!=nil{
			return err
		}
		return json.Unmarshal(data, dest)
	}
}

//ToAnyE 转换为redis支持的类型
func ToAnyE(val any) (any, error){
	var raw any
	switch v := val.(type) {
	case string:
		raw = v
	case []byte:
		raw = v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		raw = val
	default:
		// struct, map, slice 需要序列化
		buf, err := json.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("marshal fail: %w", err)
		}
		raw = buf
	}
	return raw, nil
}

//SetCmdable 设置
type SetCmdable interface{
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

func SetAny(rdb SetCmdable, ctx context.Context, key string, val any, expire time.Duration) error {
	data, err := ToAnyE(val)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, data, expire).Err()
}
