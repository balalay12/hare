package hare

import (
	"sync"
)

type queueCollection struct {
	items map[string]*queue

	mu       *sync.RWMutex
	switchMu *sync.Mutex
}

func (e *queueCollection) GetOrSet(name string, handler func() *queue) *queue {
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

func newQueueCollection() *queueCollection {
	return &queueCollection{
		items:    make(map[string]*queue),
		mu:       &sync.RWMutex{},
		switchMu: &sync.Mutex{},
	}
}
