package events

type Event interface {
	Consume(topic string, goroutine_id uint64)
}
