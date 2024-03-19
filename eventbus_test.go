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
	market := CreateMarket("test_topic", 1024)
	market.CreateConsumer(PrintEvent, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := market.PublishWithWait(ctx, "hello in")

	fmt.Println(out, err)

	market.Stop()

}
