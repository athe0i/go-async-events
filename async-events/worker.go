package async_events

import (
	"context"
	"fmt"
	"time"
)

type Worker interface {
	Run(ctx context.Context)
	HandlesErrors
}

type EventWorker struct {
	ttl   time.Duration
	poll  time.Duration
	lstnr Listener

	errChan chan error
}

func NewEventWorker(ttl time.Duration, poll time.Duration, lstnr Listener) Worker {
	ew := &EventWorker{
		ttl:   ttl,
		poll:  poll,
		lstnr: lstnr,

		errChan: make(chan error, 1),
	}

	ew.lstnr.SetErrorHandler(ew.listenerErrorHandler)

	return ew
}

func (ew *EventWorker) Run(ctx context.Context) {
	inCtx := ctx
	var cancel context.CancelFunc

	if ew.ttl > 0 {
		inCtx, cancel = context.WithTimeout(ctx, ew.ttl)
		defer cancel()
	}

	go ew.runListenLoop(inCtx)

	select {
	case <-inCtx.Done():
		return
	}
}

func (ew *EventWorker) runListenLoop(ctx context.Context) {
	for {
		ew.lstnr.Run(ctx)

		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(ew.poll)
			continue
		}
	}
}

func (ew *EventWorker) handleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-ew.errChan:
			ew.HandleError(err)
		}
	}
}

func (ew *EventWorker) listenerErrorHandler(err error) {
	ew.errChan <- err
}

func (ew *EventWorker) HandleError(err error) {
	fmt.Println(err)
}
