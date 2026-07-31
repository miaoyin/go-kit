package httputil

import (
	"io"
	"net/http"
	"time"
)


func NewTimeoutClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{Timeout: timeout},
	}
}


//DoRequest 执行请求
func DoRequest(client *http.Client, method string, url string, options ...Option) *Response {
	return CreateResponse(func() (*http.Response, error) {
		opts := NewDefaultOptions()
		for _, opt := range options {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
		req, err := opts.CreateRequest(method, url)
		if err!=nil{
			return nil, err
		}
		return client.Do(req)
	})
}

type Client struct{
	*http.Client
}
func (c *Client) Get(url string, options ...Option) *Response {
	return DoRequest(c.Client, http.MethodGet, url, options...)
}

func (c *Client) Head(url string, options ...Option) *Response {
	return DoRequest(c.Client, http.MethodHead, url, options...)
}

func (c *Client) Post(url string, contentType string, body io.Reader, options ...Option) *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return DoRequest(c.Client, http.MethodPost, url, options...)
}

func (c *Client) Put(url string, contentType string, body io.Reader, options ...Option)  *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return DoRequest(c.Client, http.MethodPut, url, options...)
}

func (c *Client) Delete(url string, contentType string, body io.Reader, options ...Option)  *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return DoRequest(c.Client, http.MethodDelete, url, options...)
}

func (c *Client) Patch(url string, contentType string, body io.Reader, options ...Option) *Response {
	options = append(options, WithContentType(contentType), WithBody(body))
	return DoRequest(c.Client, http.MethodPatch, url, options...)
}

func (c *Client) JsonPost(url string, v any, options ...Option) *Response {
	options = append(options, WithJsonBody(v))
	return DoRequest(c.Client, http.MethodPost, url, options...)
}

func (c *Client) JsonPut(url string, v any, options ...Option)  *Response {
	options = append(options, WithJsonBody(v))
	return DoRequest(c.Client, http.MethodPut, url, options...)
}

func (c *Client) JsonDelete(url string, v any, options ...Option)  *Response {
	options = append(options, WithJsonBody(v))
	return DoRequest(c.Client, http.MethodDelete, url, options...)
}

func (c *Client) JsonPatch(url string, v any, options ...Option)  *Response {
	options = append(options, WithJsonBody(v))
	return DoRequest(c.Client, http.MethodPatch, url, options...)
}
