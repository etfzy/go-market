package gomarket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/etfzy/go-market/events"
	"github.com/etfzy/go-market/pipe"
)

func echo(ctx context.Context, topic string, goroutine_id uint64, input any) (any, error) {
	fmt.Println(goroutine_id, input)
	time.Sleep(time.Second * 5)
	return input, nil
}

func TestWait(t *testing.T) {
	marketno := CreateMarket("test_topic", pipe.CreateSliceQueue(1024))
	marketno.StartConsumer(2)
	time.Sleep(time.Duration(3) * time.Second)

	e := events.CreateWaitEvent(context.TODO(), "hello wait 1", echo)
	marketno.PublishEvt(e)
	go func() {
		fmt.Println(e.GetOutput(context.TODO()))
	}()

	e2 := events.CreateWaitEvent(context.TODO(), "hello wait 2", echo)
	marketno.PublishEvt(e2)
	fmt.Println(e2.GetOutput(context.TODO()))

	marketno.Stop()
	time.Sleep(time.Second * 1)
}
