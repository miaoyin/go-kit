package eventbus

//不带topic的全局广播


//Event 事件(可选, 可自定义)
type Event struct{
    Topic     string
    Timestamp int64
    TraceID   string
    Data      any
}