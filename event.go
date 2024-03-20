package gomarket

import (
	"context"
	"errors"
)

type Event struct {
	ctx     context.Context
	input   any
	output  any
	notify  chan error
	needout bool
}

func (e *Event) WaitNotify(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("context done...")
	case err := <-e.notify:
		return err
	}
}

func (e *Event) GetOutput() any {
	return e.output
}

func createEvent(ctx context.Context, msg any, needout bool) *Event {
	if needout {
		return &Event{
			ctx:     ctx,
			input:   msg,
			notify:  make(chan error, 1),
			needout: needout,
		}
	} else {
		return &Event{
			ctx:     ctx,
			input:   msg,
			needout: needout,
		}
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
