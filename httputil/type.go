package httputil

import (
	"bytes"
	"io"
	"net/http"
)

const (
	ApplicationJson = "application/json"
)

type (
	// RequestFunc 请求
	RequestFunc func() (*http.Response, error)

	// Unmarshaller 解码
	Unmarshaller func (data []byte, v any) error
	// Marshaller 编码
	Marshaller func(v any) ([]byte, error)
)


// ToReaderE 类型转换
func ToReaderE(v any, marshaller Marshaller) (io.Reader, error) {
	rawData, err := marshaller(v)
	if err!=nil {
		return nil, err
	}
	return bytes.NewBuffer(rawData), nil
}

// ToReadCloserE 类型转换
func ToReadCloserE(v any, marshaller Marshaller) (io.ReadCloser, error) {
	rawData, err := marshaller(v)
	if err!=nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewBuffer(rawData)), nil
}