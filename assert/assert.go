package assert

import "fmt"

//MustTrue 必需为True
func MustTrue(cond bool, format string, args ...any) {
    if !cond {
        msg := fmt.Sprintf(format, args...)
        doFail("MustTrue: "+ msg)
    }
}

func MustFalse(cond bool, format string, args ...any) {
    if cond {
        msg := fmt.Sprintf(format, args...)
        doFail("MustFalse: " + msg)
    }
}

func MustNotNil(v any, format string, args ...any) {
    if v == nil {
        msg := fmt.Sprintf(format, args...)
        doFail("MustNotNil: " + msg)
    }
}

func MustNil(v any, format string, args ...any) {
    if v != nil {
        msg := fmt.Sprintf(format, args...)
        doFail("MustNil: " + msg)
    }
}

func MustEqual(a, b any, format string, args ...any) {
    if a != b {
        msg := fmt.Sprintf(format+", expect=%v, got=%v", append(args, b, a)...)
        doFail("MustEqual: " + msg)
    }
}

func MustNotEqual(a, b any, format string, args ...any) {
    if a == b {
        msg := fmt.Sprintf(format+", value=%v", append(args, a)...)
        doFail("MustNotEqual: " + msg)
    }
}

//MustLenGT 长度 > n
func MustLenGT[T any](slice []T, n int, format string, args ...any) {
    if !(len(slice) > n) {
        msg := fmt.Sprintf(format+", len=%d, require >%d", append(args, len(slice), n)...)
        doFail("MustLenGT: " + msg)
    }
}

//MustLenGE 长度 >= n
func MustLenGE[T any](slice []T, n int, format string, args ...any) {
    if !(len(slice) >= n) {
        msg := fmt.Sprintf(format+", len=%d, require >=%d", append(args, len(slice), n)...)
        doFail("MustLenGT: " + msg)
    }
}

// MustInRange 数值在 [min,max] 区间
func MustInRange[T uint64|uint32|uint|uint16|uint8|int8|int16|int|int32|int64|float32|float64](val, min, max T, format string, args ...any) {
    if val < min || val > max {
        msg := fmt.Sprintf(format+", val=%d, range [%d,%d]", append(args, val, min, max)...)
        doFail("MustInRange: " + msg)
    }
}

//MustCast 类型断言
func MustCast[T any](v any, format string, args ...any) T {
    realValue, ok := v.(T)
    if !ok {
        msg := fmt.Sprintf(format+", type=%T", append(args, v)...)
        doFail("MustCast: " + msg)
    }
    return realValue
}
