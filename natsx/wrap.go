package natsx

import (
	"context"
	"time"

	"github.com/miaoyin/go-kit/atomicx"
	"github.com/nats-io/nats.go"
)

// Snapshotter 只读
type Snapshotter interface {
	// GetSnapshot 临时获取*nats.Conn连接, 禁止长期持有
	GetSnapshot() (*nats.Conn, error)
}

// Updater 更新连接
type Updater interface {
	//SwapConn 热更新
	SwapConn(newNC *nats.Conn) (oldNC *nats.Conn)
}

type NatsClientWrap struct {
	holder *atomicx.PointerHolder[nats.Conn]
}

func NewNatsClientWrap(holder *atomicx.PointerHolder[nats.Conn]) *NatsClientWrap {
	return &NatsClientWrap{holder: holder}
}

// GetSnapshot 扩展接口
func (w *NatsClientWrap) GetSnapshot() (*nats.Conn, error) {
	return w.holder.GetSnapshot()
}

//-----------------------------------------
//    常用接口
//-----------------------------------------

func (w *NatsClientWrap) Publish(subj string, data []byte) error {
	nc, err := w.holder.GetSnapshot()
	if err != nil {
		return err
	}
	return nc.Publish(subj, data)
}

func (w *NatsClientWrap) PublishMsg(msg *nats.Msg) error {
	nc, err := w.holder.GetSnapshot()
	if err != nil {
		return err
	}
	return nc.PublishMsg(msg)
}

func (w *NatsClientWrap) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	nc, err := w.holder.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return nc.Request(subj, data, timeout)
}

func (w *NatsClientWrap) RequestWithContext(ctx context.Context, subj string, data []byte) (*nats.Msg, error) {
	nc, err := w.holder.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return nc.RequestWithContext(ctx, subj, data)
}

func (w *NatsClientWrap) Status() nats.Status {
	nc, err := w.holder.GetSnapshot()
	if err != nil {
		return nats.DISCONNECTED
	}
	return nc.Status()
}

func (w *NatsClientWrap) IsClosed() bool {
	nc, err := w.holder.GetSnapshot()
	if err != nil {
		return true
	}
	return nc.IsClosed()
}

//-----------------------------------------
//    补充方法
//-----------------------------------------

// Flush 刷新
// 扩展接口的样例, 如果NatsClientWrap提供的接口不够, 用这个方法扩展
func Flush(s Snapshotter) error {
	nc, err := s.GetSnapshot()
	if err != nil {
		return err
	}
	return nc.Flush()
}
