package httputil

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
)

type (
    Option func(p *Options) error
)

func NewDefaultOptions() Options {
    return Options{
        Context: context.Background(),
        Body: nil,
    }
}


//Options 请求参数选项
type Options struct{
    //请求context
    Context     context.Context
    //请求Header
    Header      http.Header
    //请求内容
    Body        io.Reader
    //请求内容格式
    ContentType string

}

func (p Options) CreateRequest(method string, url string) (*http.Request, error){
    req, err := http.NewRequestWithContext(p.Context,method, url, p.Body)
    if err!=nil{
        return nil, err
    }
    return req, nil
}



// WithHeader 头部信息
func WithHeader(header http.Header) Option {
    return func(op *Options) error{
        op.Header = header
        return nil
    }
}

// WithContentType 内容格式
func WithContentType(contentType string) Option {
    return func(op *Options) error{
        if len(contentType)>0{
            op.Header.Set("Content-Type", contentType)
        }
        return nil
    }
}

// WithBody 请求body
func WithBody(body io.Reader) Option {
    return func(op *Options) error{
        op.Body = body
        return nil
    }
}

// WithJsonBody 请求参数json编码
func WithJsonBody(v any) Option {
    return func(op *Options) error{
        op.Header.Add("Content-Type", ApplicationJson)
        body, err := ToReaderE(v, json.Marshal)
        if err != nil {
            return err
        }
        op.Body = body
        return nil
    }
}

// WithMarshallerBody 请求参数编码
func WithMarshallerBody(v any, marshaller Marshaller) Option {
    return func(op *Options) error{
        body, err := ToReaderE(v, marshaller)
        if err != nil {
            return err
        }
        op.Body = body
        return nil
    }
}


// WithContext 请求上下文
func WithContext(ctx context.Context) Option {
    return func(op *Options) error{
        op.Context = ctx
        return nil
    }
}
