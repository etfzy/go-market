package gomarket

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func PrintEvent(ctx context.Context, goroutine_id uint64, input any) (any, error) {
	fmt.Println("consumer receive msg...", goroutine_id)
	in := input.(string)
	fmt.Println(in)

	return "PrintEvent hello out", nil
}

func TestEvent(t *testing.T) {
	market := CreateMarket("test_topic", 1024, PrintEvent)
	market.StartConsumer(3)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := market.PublishWithWait(ctx, "hello in")

	fmt.Println(out, err)

	fmt.Println("adjust")
	market.AdjustConsumer(1)

	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("stop")
	market.Stop()

	time.Sleep(time.Duration(10) * time.Second)

}
