package must

import (
	"errors"
	"fmt"
	"io"
)

// 减少if, 适合单元测

// Get 减少if err判断
func Get[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Errorf("must: %w", err))
	}
	return v
}

// OK 减少if ok判断
func OK[T any](v T, ok bool) T {
	if !ok {
		panic(errors.New("must: value not ok"))
	}
	return v
}

// NoErr 减少 if err判断
func NoErr(err error) {
	if err != nil {
		panic(fmt.Errorf("must: %w", err))
	}
}

func True(cond bool, msg ...any) {
	if !cond {
		panic(fmt.Errorf("must: condition failed: %v", msg))
	}
}

func False(cond bool, msg ...any) {
	if cond {
		panic(fmt.Errorf("must: condition should be false: %v", msg))
	}
}

func NotZero[T comparable](v T, msg ...any) {
	var zero T
	if v == zero {
		panic(fmt.Errorf("must: zero value, %v", msg))
	}
}

func NotNil[T any](v *T, msg ...any) *T {
	if v == nil {
		panic(fmt.Errorf("must: nil pointer, %v", msg))
	}
	return v
}

func Equal[T comparable](a, b T, msg ...any) {
	if a != b {
		panic(fmt.Errorf("must: %v != %v, %v", a, b, msg))
	}
}

func Close(c io.Closer) {
	NoErr(c.Close())
}

func Panicf(format string, args ...any) {
	panic(fmt.Errorf("must: "+format, args...))
}
