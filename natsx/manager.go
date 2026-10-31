package natsx

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/miaoyin/go-kit/atomicx"
	"github.com/miaoyin/go-kit/lifecycle"
	"github.com/nats-io/nats.go"
)

type NatsConfig struct {
	URL      string
	User     string
	Password string
	Timeout  time.Duration
}

type OnConnChanged func(nc *nats.Conn)

// Manager 管理连接
type Manager struct {
	*lifecycle.LifeCycler
	holder    *atomicx.Holder[nats.Conn]
	reloadCh  chan struct{}
	curConf   atomic.Pointer[NatsConfig]
	onChanged OnConnChanged
	runCancel context.CancelFunc
}

func NewManager(initCfg *NatsConfig, holder *atomicx.Holder[nats.Conn]) *Manager {
	m := &Manager{
		holder:   holder,
		reloadCh: make(chan struct{}, 1),
	}
	m.LifeCycler = lifecycle.NewLifeCycler(lifecycle.Hook{
		OnStart: m.onStart,
		OnClose: m.onClose,
	})
	m.curConf.Store(initCfg)
	return m
}

func (m *Manager) SetOnChanged(cb OnConnChanged) {
	m.onChanged = cb
}

func (m *Manager) Reload() {
	select {
	case m.reloadCh <- struct{}{}:
	default:
	}
}

func (m *Manager) UpdateConfig(cfg *NatsConfig) {
	m.curConf.Store(cfg)
	m.Reload()
}

func (m *Manager) onStart(_ context.Context) error {
	runCtx, runCancel := context.WithCancel(context.Background())
	m.runCancel = runCancel
	go m.runLoop(runCtx)
	m.Reload()
	return nil
}

func (m *Manager) runLoop(runCtx context.Context) {
	for {
		select {
		case <-runCtx.Done():
			return
		case <-m.reloadCh:
			if err := m.doReload(); err != nil {
				fmt.Printf("nats reload fail: %v\n", err)
			}
		}
	}
}

func (m *Manager) doReload() error {
	newConf := m.curConf.Load()
	if newConf == nil {
		return errors.New("empty config")
	}
	opts := []nats.Option{nats.Timeout(newConf.Timeout)}
	if newConf.User != "" && newConf.Password != "" {
		opts = append(opts, nats.UserInfo(newConf.User, newConf.Password))
	}
	newNC, err := nats.Connect(newConf.URL, opts...)
	if err != nil {
		return err
	}
	if err = newNC.FlushTimeout(2 * time.Second); err != nil {
		newNC.Close()
		return fmt.Errorf("nats flush check fail: %w", err)
	}

	oldNC := m.holder.Swap(newNC)

	if m.onChanged != nil {
		m.onChanged(newNC)
	}

	if oldNC != nil {
		go func(old *nats.Conn) {
			time.Sleep(3 * time.Second)
			old.Close()
		}(oldNC)
	}
	return nil
}

func (m *Manager) onClose(ctx context.Context) error {
	m.runCancel()
	nc, err := m.holder.GetSnapshot()
	if err == nil && nc != nil {
		_ = nc.Drain()
		nc.Close()
	}
	return nil
}
