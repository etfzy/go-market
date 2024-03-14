package gomarket

import (
	"context"
	"errors"
)

type Market struct {
	topic  string
	queues chan *Event
	cap    uint64
}

func CreateMarket(topic string, cap uint64) *Market {
	if cap < 1 {
		panic("cap must >= 1")
	}
	return &Market{
		topic:  topic,
		queues: make(chan *Event, cap),
	}
}

func (m *Market) CreateConsumer(fn ConsumerFn, number int64) {
	if number < 1 {
		panic("consumer number must >= 1")
	}

	for i := 0; i < int(number); i++ {
		c := CreateConsumer(m.topic, uint64(i), fn, m.queues)
		go c.ConsumerRun()
	}
}

func (m *Market) PublishWithWait(ctx context.Context, msg any) (any, error) {
	event := createEvent(ctx, msg)

	err := ctx.Err()

	if err != nil {
		return nil, err
	}

	select {
	case m.queues <- event:
	default:
		return nil, errors.New("market queue is full...")
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("context done ")
	case err := <-event.notify:
		return event.output, err
	}
}
