package pipe

import (
	"errors"

	"github.com/etfzy/go-market/events"
)

type SliceQueue struct {
	queue []events.Event
	cap   int
}

func CreateSliceQueue(cap int) MsgQueue {
	return &SliceQueue{
		queue: make([]events.Event, 0, cap),
		cap:   cap,
	}
}

func (sq *SliceQueue) PublishEvt(e events.Event) error {
	if len(sq.queue) >= sq.cap {
		return errors.New("exceed max cap...")
	}

	sq.queue = append(sq.queue, e)
	return nil
}

func (sq *SliceQueue) CancelEvt(e events.Event) error {
	for _, v := range sq.queue {
		if v == e {
			return nil
		}
	}
	return errors.New("not find...")

}

func (sq *SliceQueue) PopEvt() events.Event {
	if len(sq.queue) == 0 {
		return nil
	}

	temp := sq.queue[0]
	sq.queue = sq.queue[1:]
	return temp
}
