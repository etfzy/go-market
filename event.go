package gomarket

import (
	"context"
	"errors"
)

type Event struct {
	ctx    context.Context
	input  any
	output any
	notify chan error
}

func (e *Event) WaitNotify(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("context done...")
	case err := <-e.notify:
		return err
	}
}

func createEvent(ctx context.Context, msg any) *Event {
	return &Event{
		ctx:    ctx,
		input:  msg,
		notify: make(chan error, 1),
	}
}

func (e *Event) reponse(output any, err error) {
	e.output = output
	select {
	case e.notify <- err:
	default:
		return
	}
}
