package gomarket

import (
	"fmt"
	"testing"
	"time"

	"github.com/etfzy/go-market/pipe"
)

type PriorityEvt struct {
	priority int
}

func (ne *PriorityEvt) Consume(topic string, goroutine_id uint64) {
	fmt.Println("---------", topic, goroutine_id, ne.priority, "---------")

}

func (ne *PriorityEvt) GetPriority() int {
	return ne.priority
}

func TestPriorityNowait(t *testing.T) {
	fmt.Println("test priority")
	marketno := CreateMarket("test_topic", pipe.CreatePriorityQueue(1024))
	marketno.StartConsumer(1)

	pe3 := &PriorityEvt{
		priority: 3,
	}

	pe1 := &PriorityEvt{
		priority: 1,
	}

	pe2 := &PriorityEvt{
		priority: 2,
	}

	marketno.PublishEvt(pe2)
	marketno.PublishEvt(pe1)
	marketno.PublishEvt(pe3)
	marketno.CancelEvt(pe1)

	time.Sleep(1*time.Second)
}
