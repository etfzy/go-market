package pipe

import (
	"sync"

	"github.com/etfzy/go-market/events"
)

type MsgQueue interface {
	PublishEvt(e events.Event) error
	CancelEvt(e events.Event) error
	PopEvt() events.Event
}

type Pipe struct {
	CondLock *sync.Mutex
	Cond     *sync.Cond

	StateLock *sync.RWMutex
	Running   bool

	QueueLock *sync.Mutex
	Queue     MsgQueue
}

func CreatePipe(q MsgQueue) *Pipe {
	condlock := &sync.Mutex{}
	return &Pipe{
		CondLock:  condlock,
		Cond:      sync.NewCond(condlock),
		StateLock: &sync.RWMutex{},
		Running:   true,
		QueueLock: &sync.Mutex{},
		Queue:     q,
	}
}

func (p *Pipe) CheckState() bool {
	p.StateLock.RLock()
	defer p.StateLock.RUnlock()
	return p.Running
}

func (p *Pipe) WaitCond() {
	p.CondLock.Lock()
	defer p.CondLock.Unlock()
	p.Cond.Wait()
}

func (p *Pipe) Stop() {
	p.StateLock.Lock()
	defer p.StateLock.Unlock()
	p.Running = false
}

func (p *Pipe) Broadcast() {
	p.Cond.Broadcast()
	return
}

func (p *Pipe) publishQueueEvt(evt events.Event) error {
	p.QueueLock.Lock()
	defer p.QueueLock.Unlock()
	return p.Queue.PublishEvt(evt)
}

func (p *Pipe) PublishEvt(evt events.Event) error {
	err := p.publishQueueEvt(evt)
	if err != nil {
		return err
	} else {
		p.Cond.Signal()
		return nil
	}
}

func (p *Pipe) CancelEvt(evt events.Event) error {
	p.QueueLock.Lock()
	defer p.QueueLock.Unlock()
	return p.Queue.CancelEvt(evt)
}

// consumer 外层已经加过锁
func (p *Pipe) PopEvt() events.Event {
	p.QueueLock.Lock()
	defer p.QueueLock.Unlock()
	return p.Queue.PopEvt()
}
