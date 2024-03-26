package pipe

import (
	"container/heap"
	"errors"

	"github.com/etfzy/go-market/events"
)

type PrioriEvent interface {
	GetPriority() int
}

// Item represents an item in the priority queue.
type Item struct {
	evt      events.Event // 值
	priority int          // 优先级
}

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue struct {
	cap   int
	items []Item
}

// Len returns the length of the priority queue.
func (pq *PriorityQueue) Len() int { return len(pq.items) }

// Less returns true if the item with index i has higher priority than the item with index j.
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].priority > pq.items[j].priority
}

// Swap swaps the items with indexes i and j in the priority queue.
func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

// Push pushes the given item onto the priority queue.
func (pq *PriorityQueue) Push(x any) {
	pq.items = append(pq.items, x.(Item))
}

func (pq *PriorityQueue) Remove(evt events.Event) bool {

	for k, v := range pq.items {
		if v.evt == evt {
			pq.items = append(pq.items[:k],pq.items[k+1:]...)
			return true
		}
	}
	return false
}

// Pop removes and returns the highest priority item from the priority queue.
func (pq *PriorityQueue) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

func CreatePriorityQueue(cap int) *PriorityQueue {

	pq := &PriorityQueue{
		cap:   cap,
		items: make([]Item, 0, cap),
	}

	return pq
}

func (pq *PriorityQueue) PublishEvt(e events.Event) error {
	if len(pq.items) >= pq.cap {
		return errors.New("exceed max cap...")
	}

	item := Item{
		evt:      e,
		priority: e.(PrioriEvent).GetPriority(),
	}

	heap.Push(pq, item)
	return nil
}

func (pq *PriorityQueue) CancelEvt(e events.Event) error {
	if pq.Remove(e) {
		return nil
	}

	return errors.New("not find...")

}

func (pq *PriorityQueue) PopEvt() events.Event {
	if len(pq.items) == 0 {
		return nil
	}

	temp := heap.Pop(pq)

	if temp == nil {
		return nil
	}

	item := temp.(Item)
	return item.evt
}
