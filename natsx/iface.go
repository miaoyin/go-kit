package natsx

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher interface {
	Publish(subj string, data []byte) error
	PublishMsg(msg *nats.Msg)
}

type Requester interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	RequestWithContext(ctx context.Context, subj string, data []byte) (*nats.Msg, error)
}

type Subscriber interface {
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error)
}

type JetStreamPusher interface {
	PublishAsync(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAckFuture, error)
}
