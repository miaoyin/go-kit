package chutil

import (
    "context"
    "errors"
    "time"
)

var (
    ErrChannelClosed = errors.New("channel closed")
)

//WaitChan 等待channel(context+channel+timeout)
func WaitChan[T any](ctx context.Context, ch <-chan T, timeout time.Duration) (val T, err error) {
    timer := time.NewTimer(timeout)
    defer timer.Stop()

    select {
    case <-ctx.Done():
        return val, ctx.Err()
    case <-timer.C:
        return val, context.DeadlineExceeded
    case v, ok := <-ch:
        if !ok {
            return val, ErrChannelClosed
        }
        return v, nil
    }
}


//SleepOrCancel 等待(context+timeout)
// 带context的Sleep
func SleepOrCancel(ctx context.Context, d time.Duration) error {
    timer := time.NewTimer(d)
    defer timer.Stop()

    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-timer.C:
        return nil
    }
}
