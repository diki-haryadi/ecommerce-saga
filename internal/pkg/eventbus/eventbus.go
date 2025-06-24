package eventbus

import (
	"sync"
)

type EventBus struct {
	subscribers map[string][]chan interface{}
	mu          sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan interface{}),
	}
}

func (b *EventBus) Subscribe(topic string) chan interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan interface{}, 1)
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	return ch
}

func (b *EventBus) Publish(topic string, data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subscribers, exists := b.subscribers[topic]; exists {
		for _, ch := range subscribers {
			go func(ch chan interface{}) {
				ch <- data
			}(ch)
		}
	}
}
