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
	time.Sleep(3 * time.Second)
}

func TestNoWait(t *testing.T) {
	marketno := CreateMarket("test_topic", pipe.CreateSliceQueue(1024))
	marketno.StartConsumer(2)
	time.Sleep(time.Duration(3) * time.Second)

	e := &NoWaitEvt{
		in: "hello no wiat 1",
	}
	marketno.PublishEvt(e)
	fmt.Println("first publish")

	e2 := &NoWaitEvt{
		in: "hello no wiat 2",
	}
	marketno.PublishEvt(e2)
	fmt.Println("second publish")

	time.Sleep(1 * time.Second)

	marketno.Stop()

	time.Sleep(1 * time.Second)

}
