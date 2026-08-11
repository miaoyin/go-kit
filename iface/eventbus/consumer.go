package eventbus

import (
	"context"

	"github.com/miaoyin/go-kit/iface/logger"
)

type EventConsumerConfiguration[T any] struct{
	//线程数
	WorkerNum int
	//队列大小
	QueueSize int
	//日志(可选)
	Logger logger.Logger
	//事件处理函数
	OnEvent func(evt *T)
}

func NewEventConsumerConfiguration[T any]() EventConsumerConfiguration[T]{
	return EventConsumerConfiguration[T]{
		WorkerNum: 1,
		QueueSize: 256,
		Logger: &logger.NopLogger{},
	}
}

//EventConsumer 样例
type EventConsumer[T any] struct{
	queue chan *T
	config EventConsumerConfiguration[T]
	ctx   context.Context
	cancel context.CancelFunc
}

//NewEventConsumer
//  workerNum 线程数量
//  queueSize 队列大小
func NewEventConsumer[T any](config EventConsumerConfiguration[T]) *EventConsumer[T] {
	ctx,cancel := context.WithCancel(context.Background())
	l := &EventConsumer[T]{
		queue: make(chan *T, config.QueueSize),
		ctx: ctx,
		cancel: cancel,
	}
	for i:=0;i<config.WorkerNum;i++ {
		go l.worker()
	}
	return l
}

//Deliver 投入事件
func (l *EventConsumer[T]) Deliver(evt *T) {
	select {
	case l.queue <- evt:
	default:
		l.config.Logger.Warn("queue is full", "evt", evt)
	}
}

func (l *EventConsumer[T]) worker() {
	for {
		select {
		case <- l.ctx.Done():
			return
		case evt := <- l.queue:
			l.config.OnEvent(evt)
		}
	}
}

func (l *EventConsumer[T]) Stop() {
	l.cancel()
}