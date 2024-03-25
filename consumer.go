package gomarket

import (
	"fmt"

	"github.com/etfzy/go-market/pipe"
)

type Consumer struct {
	idx   uint64
	topic string
	stop  chan int
	pipe  *pipe.Pipe
}

func CreateConsumer(topic string, idx uint64, pipe *pipe.Pipe) *Consumer {
	return &Consumer{
		topic: topic,
		idx:   idx,
		stop:  make(chan int, 1),
		pipe:  pipe,
	}
}

func (c *Consumer) ConsumerRun() {

	for c.pipe.CheckState() {

		for {
			if !c.pipe.CheckState() {
				break
			}

			event := c.pipe.PopEvt()
			if event != nil {
				event.Consume(c.topic, c.idx)
			} else {
				break
			}
		}

		if !c.pipe.CheckState() {
			break
		}

		c.pipe.WaitCond()

		if !c.pipe.CheckState() {
			break
		}
	}

	fmt.Println("consumer exit ", c.topic, c.idx)
}
