package httpclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// Doer 执行请求
type Doer interface {
	Do(ctx context.Context, method, reqURL string, body io.Reader, headers map[string]string) (*Response, error)
}

type HTTPClient interface {
	Do(ctx context.Context, method, reqURL string, body io.Reader, headers map[string]string) (*Response, error)

	Get(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error)
	Head(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error)

	PostJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error)
	PostForm(ctx context.Context, reqURL string, form map[string]string, headers map[string]string) (*Response, error)
	PutJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error)
	PatchJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error)
	Delete(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error)
}

//----------------------------------------
//   返回结果
//----------------------------------------

type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func (r *Response) BindJSON(target any) error {
	return json.Unmarshal(r.Body, target)
}
