package hare

import (
	"sync"
)

type consumerCollection struct {
	items map[string]*consumer

	mu       *sync.RWMutex
	switchMu *sync.Mutex
}

func (e *consumerCollection) GetOrSet(name string, handler func() *consumer) *consumer {
	e.mu.RLock()
	if instance, ok := e.items[name]; ok {
		e.mu.RUnlock()
		return instance
	}
	{
		e.switchMu.Lock()
		e.mu.RUnlock()
		e.mu.Lock()
		defer e.mu.Unlock()
		e.switchMu.Unlock()
	}
	e.items[name] = handler()
	return e.items[name]
}

func newConsumerCollection() *consumerCollection {
	return &consumerCollection{
		items:    make(map[string]*consumer),
		mu:       &sync.RWMutex{},
		switchMu: &sync.Mutex{},
	}
}
