package gomarket

import (
	"fmt"
	"testing"
	"time"

	"github.com/etfzy/go-market/pipe"
)

type NoWaitEvt struct {
	in string
}

func (ne *NoWaitEvt) Consume(topic string, goroutine_id uint64) {
	fmt.Println("---------", topic, goroutine_id, ne.in, "---------")
}

func TestNoWait(t *testing.T) {
	marketno := CreateMarket("test_topic", pipe.CreateSliceQueue(1024))
	marketno.StartConsumer(2)

	e := &NoWaitEvt{
		in: "hello no wiat 1",
	}
	marketno.PublishEvt(e)

	e2 := &NoWaitEvt{
		in: "hello no wiat 2",
	}
	marketno.PublishEvt(e2)

	e3 := &NoWaitEvt{
		in: "hello no wiat 3",
	}
	marketno.PublishEvt(e3)

	marketno.CancelEvt(e2)

	time.Sleep(1*time.Second)

}
