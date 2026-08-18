package service

import (
	"encoding/json"
	"sync"
)

type AppEvent struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type EventBus struct {
	mu   sync.RWMutex
	next int
	subs map[int]chan AppEvent
}

func NewEventBus() *EventBus { return &EventBus{subs: map[int]chan AppEvent{}} }
func (b *EventBus) Subscribe() (int, <-chan AppEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.next
	b.next++
	ch := make(chan AppEvent, 32)
	b.subs[id] = ch
	return id, ch
}
func (b *EventBus) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if ch, ok := b.subs[id]; ok {
		delete(b.subs, id)
		close(ch)
	}
}
func (b *EventBus) Publish(event AppEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs {
		select {
		case ch <- event:
		default:
		}
	}
}
func (e AppEvent) JSON() []byte { b, _ := json.Marshal(e.Data); return b }
