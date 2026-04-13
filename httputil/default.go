package httputil

import (
    "io"
    "net/http"
    "time"
)

var (
    // DefaultClient http客户端
    DefaultClient = &Client{
        Client: &http.Client{},
    }
    DefaultTimeoutClient = NewTimeoutClient(time.Second * 30)
)

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
