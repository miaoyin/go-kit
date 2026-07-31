package chutil

import "context"

type (
	Item struct {
		V interface{}
		E error
	}
)

func Of(i interface{}) Item {
	return Item{V: i}
}

func Error(err error) Item {
	return Item{E: err}
}

func (i Item) Error() bool {
	return i.E != nil
}

func (i Item) SendBlocking(ch chan<- Item) {
	ch <- i
}

func (i Item) SendContext(ctx context.Context, ch chan<- Item) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		select {
		case <-ctx.Done():
			return false
		case ch <- i:
			return true
		}
	}
}

func (i Item) SendNonBlocking(ch chan<- Item) bool {
	select {
	default:
		return false
	case ch <- i:
		return true
	}
}
