package gomarket

import (
	"context"
	"fmt"
)

type ConsumerFn func(ctx context.Context, goroutine_id uint64, input any) (any, error)

type Consumer struct {
	idx   uint64
	topic string
	fn    ConsumerFn
	queue chan *Event
	stop  chan int
}

func CreateConsumer(topic string, idx uint64, fn ConsumerFn, queue chan *Event) *Consumer {
	return &Consumer{
		topic: topic,
		idx:   idx,
		fn:    fn,
		queue: queue,
		stop:  make(chan int, 1),
	}
}

func (c *Consumer) Stop() {
	fmt.Println("prepare consumer stop ", c.topic, c.idx)

	select {
	case c.stop <- 1:
	default:
	}
}

func (c *Consumer) ConsumerRun() {
	fmt.Println("consumer start ", c.topic, c.idx)
loop:
	for {
		select {
		case event := <-c.queue:
			out, err := c.fn(event.ctx, c.idx, event.input)
			event.reponse(out, err)
		case <-c.stop:
			break loop
		}
	}

	fmt.Println("consumer exit ", c.topic, c.idx)
}
