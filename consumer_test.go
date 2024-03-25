package gomarket

import (
	"fmt"
	"testing"
	"time"

	"github.com/etfzy/go-market/pipe"
)

func TestConsumer(t *testing.T) {
	market := CreateMarket("test_topic", pipe.CreateSliceQueue(1024))
	market.StartConsumer(3)
	time.Sleep(time.Duration(3) * time.Second)
	market.pipe.Broadcast()

	time.Sleep(time.Duration(3) * time.Second)
	fmt.Println("stop")
	market.Stop()
}
