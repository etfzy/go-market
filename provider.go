package gomarket

import (
	"github.com/etfzy/go-market/events"
	"github.com/etfzy/go-market/pipe"
)

type Market struct {
	topic string

	consumers []*Consumer

	pipe *pipe.Pipe
}

func CreateMarket(topic string, queue pipe.MsgQueue) *Market {
	market := &Market{
		topic:     topic,
		consumers: make([]*Consumer, 0, 1024),
		pipe:      pipe.CreatePipe(queue),
	}

	return market
}

func (m *Market) StartConsumer(number int64) {
	if number < 1 {
		panic("consumer number must >= 1")
	}

	for i := 0; i < int(number); i++ {
		c := CreateConsumer(m.topic, uint64(i), m.pipe)
		m.consumers = append(m.consumers, c)
		go c.ConsumerRun()
	}
}

func (m *Market) Stop() {
	m.pipe.Stop()
	m.pipe.Broadcast()
	m.consumers = m.consumers[:0]
}

func (m *Market) PublishEvt(evt events.Event) error {
	return m.pipe.PublishEvt(evt)
}

func (m *Market) CancelEvt(evt events.Event) error {
	return m.pipe.CancelEvt(evt)
}
