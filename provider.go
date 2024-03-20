package gomarket

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Market struct {
	topic  string
	queues chan *Event

	consumber_lock sync.Mutex
	consumers      []*Consumer
	cap            uint64
	fn             ConsumerFn
}

func CreateMarket(topic string, cap uint64, fn ConsumerFn) *Market {
	if cap < 1 {
		panic("cap must >= 1")
	}
	return &Market{
		topic:          topic,
		queues:         make(chan *Event, cap),
		fn:             fn,
		consumers:      make([]*Consumer, 0, 1024),
		consumber_lock: sync.Mutex{},
	}
}

func (m *Market) StartConsumer(number int64) {
	if number < 1 {
		panic("consumer number must >= 1")
	}

	m.consumber_lock.Lock()
	defer m.consumber_lock.Unlock()
	for i := 0; i < int(number); i++ {
		c := CreateConsumer(m.topic, uint64(i), m.fn, m.queues)
		m.consumers = append(m.consumers, c)
		go c.ConsumerRun()
	}
}

func (m *Market) AdjustConsumer(expect int64) {
	if expect == 0 {
		m.Stop()
		return
	}

	m.consumber_lock.Lock()
	defer m.consumber_lock.Unlock()
	current := int64(len(m.consumers))

	if current == expect {
		return
	} else if expect > current {
		add := expect - current
		for i := 0; i < int(add); i++ {
			c := CreateConsumer(m.topic, uint64(i), m.fn, m.queues)
			m.consumers = append(m.consumers, c)
			go c.ConsumerRun()
		}
	} else {
		del := current - expect

		del_consumers := m.consumers[:del]

		for _, c := range del_consumers {
			c.Stop()
		}
		m.consumers = m.consumers[del:]
	}

	fmt.Println("current consumer number", len(m.consumers))

}

func (m *Market) Stop() {
	m.consumber_lock.Lock()
	defer m.consumber_lock.Unlock()
	for _, c := range m.consumers {
		c.Stop()
	}

	m.consumers = m.consumers[:0]

}

func (m *Market) PublishWithWait(ctx context.Context, msg any) (any, error) {
	event := createEvent(ctx, msg, true)

	err := ctx.Err()

	if err != nil {
		return nil, err
	}

	select {
	case m.queues <- event:
	default:
		return nil, errors.New("market queue is full...")
	}

	err = event.WaitNotify(ctx)

	return event.output, err
}

func (m *Market) PublishAsync(ctx context.Context, msg any) (*Event, error) {
	event := createEvent(ctx, msg, true)

	err := ctx.Err()

	if err != nil {
		return nil, err
	}

	select {
	case m.queues <- event:
	default:
		return nil, errors.New("market queue is full...")
	}

	return event, nil
}

func (m *Market) PublishWithoutOutput(ctx context.Context, msg any) error {
	event := createEvent(ctx, msg, false)

	err := ctx.Err()

	if err != nil {
		return err
	}

	select {
	case m.queues <- event:
	default:
		return errors.New("market queue is full...")
	}

	return nil
}
