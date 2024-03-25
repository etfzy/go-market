package events

import (
	"context"
	"errors"
)

type CustomWaitFn func(ctx context.Context, topic string, goroutine_id uint64, input any) (any, error)

type WaitEvent struct {
	ctx      context.Context
	input    any
	output   any
	notify   chan error
	customFn CustomWaitFn
}

func (e *WaitEvent) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("context done...")
	case err := <-e.notify:
		return err
	}
}

func (e *WaitEvent) GetOutput(ctx context.Context) (any, error) {
	err := e.wait(ctx)

	if err != nil {
		return nil, err
	} else {
		return e.output, nil
	}
}

func CreateWaitEvent(ctx context.Context, msg any, cf CustomWaitFn) *WaitEvent {
	return &WaitEvent{
		ctx:      ctx,
		input:    msg,
		notify:   make(chan error, 1),
		customFn: cf,
	}
}

func (e *WaitEvent) Consume(topic string, goroutine_id uint64) {
	out, err := e.customFn(e.ctx, topic, goroutine_id, e.input)
	e.reponse(out, err)
}

func (e *WaitEvent) reponse(output any, err error) {
	e.output = output
	select {
	case e.notify <- err:
	default:
		return
	}
}
