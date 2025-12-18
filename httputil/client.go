package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)


var (
	// DefaultClient http客户端
	DefaultClient = &Client{
		Client: &http.Client{},
	}
	DefaultTimeoutClient = TimeoutClient(time.Second * 30)
)

func TimeoutClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{Timeout: timeout},
	}
}


type (
	Option func(req *Request) error
)

// WithHeader 头部信息
func WithHeader(header http.Header) Option {
	return func(req *Request) error{
		req.Header = header
		return nil
	}
}

// WithContentType 内容格式
func WithContentType(contentType string) Option {
	return func(req *Request) error{
		if len(contentType)>0{
			req.Header.Set("Content-Type", contentType)
		}
		return nil
	}
}

// WithBody 请求body
func WithBody(body io.Reader) Option {
	return func(req *Request) error{
		req.SetBody(body)
		return nil
	}
}

// WithJsonBody 请求参数json编码
func WithJsonBody(v any) Option {
	return func(req *Request) error{
		req.Header.Add("Content-Type", ApplicationJson)
		body, err := ToReaderE(v, json.Marshal)
		if err != nil {
			return err
		}
		return WithBody(body)(req)
	}
}

// WithMarshallerBody 请求参数编码
func WithMarshallerBody(v any, marshaller Marshaller) Option {
	return func(req *Request) error{
		body, err := ToReaderE(v, marshaller)
		if err != nil {
			return err
		}
		return WithBody(body)(req)
	}
}


// WithContext 请求上下文
func WithContext(ctx context.Context) Option {
	return func(req *Request) error {
		req.Request = req.Request.WithContext(ctx)
		return nil
	}
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err!=nil{
		return nil, err
	}
	return &Request{
		Request:req,
	}, nil
}

type Request struct{
	*http.Request
}

func (r *Request) SetBody(body io.Reader) {
	if body == nil {
		return
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	r.Body = rc
	switch v := body.(type) {
	case *bytes.Buffer:
		r.ContentLength = int64(v.Len())
		buf := v.Bytes()
		r.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return io.NopCloser(r), nil
		}
	case *bytes.Reader:
		r.ContentLength = int64(v.Len())
		snapshot := *v
		r.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	case *strings.Reader:
		r.ContentLength = int64(v.Len())
		snapshot := *v
		r.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	default:
	}
	if r.GetBody != nil && r.ContentLength == 0 {
		r.Body = http.NoBody
		r.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	}

}

type Client struct{
	*http.Client
}

func (c *Client) Do(method string, url string, options ...Option) *Response{
	return CreateResponse(func() (*http.Response, error) {
		req, err := NewRequest(method, url, nil)
		if err!=nil{
			return nil, err
		}
		for _, option := range options {
			if err = option(req);err!=nil{
				return nil, err
			}
		}
		return c.Client.Do(req.Request)
	})
}

func (c *Client) Get(url string, options ...Option) *Response{
	return c.Do(http.MethodGet, url, options...)
}

func (c *Client) Head(url string, options ...Option) *Response {
	return c.Do(http.MethodHead, url, options...)
}

func (c *Client) Post(url string, contentType string, body io.Reader, options ...Option) *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return c.Do(http.MethodPost, url, options...)
}

func (c *Client) Put(url string, contentType string, body io.Reader, options ...Option)  *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return c.Do(http.MethodPut, url, options...)
}

func (c *Client) Delete(url string, contentType string, body io.Reader, options ...Option)  *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return c.Do(http.MethodDelete, url, options...)
}

func (c *Client) Patch(url string, contentType string, body io.Reader, options ...Option) *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return c.Do(http.MethodPatch, url, options...)
}

func (c *Client) JsonPost(url string, v any, options ...Option) *Response {
	options = append(options, WithJsonBody(v))
	return c.Do(http.MethodPost, url, options...)
}

func (c *Client) JsonPut(url string, v any, options ...Option)  *Response {
	options = append(options, WithJsonBody(v))
	return c.Do(http.MethodPut, url, options...)
}

func (c *Client) JsonDelete(url string, v any, options ...Option)  *Response {
	options = append(options, WithJsonBody(v))
	return c.Do(http.MethodDelete, url, options...)
}

func (c *Client) JsonPatch(url string, v any, options ...Option)  *Response {
	options = append(options, WithJsonBody(v))
	return c.Do(http.MethodPatch, url, options...)
}


// Get 请求get
func Get(url string) *Response {
	return DefaultClient.Get(url)
}

// Post 请求post
func Post(url string, contentType string, body io.Reader) *Response {
	return DefaultClient.Post(url, contentType, body)
}

// Head 请求Head
func Head(url string) *Response {
	return DefaultClient.Head(url)
}

// Put 请求Put
func Put(url string, contentType string, body io.Reader) *Response {
	return DefaultClient.Put(url, contentType, body)
}

// Delete 请求delete
func Delete(url string, contentType string, body io.Reader) *Response {
	return DefaultClient.Delete(url, contentType, body)
}

// Patch 请求Patch
func Patch(url string, contentType string, body io.Reader) *Response {
	return DefaultClient.Patch(url, contentType, body)
}

// JsonPost 请求json
func JsonPost(url string, v any) *Response {
	return DefaultClient.JsonPost(url, v)
}

// JsonPut 请求Put
func JsonPut(url string, v any) *Response {
	return DefaultClient.JsonPut(url, v)
}

// JsonDelete 请求delete
func JsonDelete(url string, v any) *Response {
	return DefaultClient.JsonDelete(url, v)
}

// JsonPatch 请求patch
func JsonPatch(url string, v any) *Response {
	return DefaultClient.JsonPatch(url, v)
}

