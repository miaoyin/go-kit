package httputil

import (
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
