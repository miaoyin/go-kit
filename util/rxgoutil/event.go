package rxgoutil

import "github.com/reactivex/rxgo/v2"


//EventObservable 适合事件广播
type EventObservable struct {
    //Src事件源
    Src chan rxgo.Item
    rxgo.Observable
}

func NewEventObservable(size int) *EventObservable{
    src := make(chan rxgo.Item, size)
    return &EventObservable{
        Src: src,
        Observable: rxgo.FromEventSource(src, rxgo.WithBackPressureStrategy(rxgo.Drop)),
    }
}

//SendBlocking 发布事件, 阻塞
//  需要设置Drop, 确保将消息发送给监听者
//  1.有监听: 确保传给监听者
//  2.无监听: 丢弃
func (ev *EventObservable) SendBlocking(v interface{})  {
    rxgo.Of(v).SendBlocking(ev.Src)
}

//SendNonBlocking 发布事件, 非阻塞
func (ev *EventObservable) SendNonBlocking(v interface{})  bool{
    return rxgo.Of(v).SendNonBlocking(ev.Src)
}

//Close 关闭
func (ev *EventObservable) Close(){
    close(ev.Src)
}
