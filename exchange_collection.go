package hare

import (
	"sync"
)

type exchangeCollection struct {
	items map[string]*exchange

	mu       *sync.RWMutex
	switchMu *sync.Mutex
}

func (e *exchangeCollection) GetOrSet(name string, handler func() *exchange) *exchange {
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

func newExchangeCollection() *exchangeCollection {
	return &exchangeCollection{
		items:    make(map[string]*exchange),
		mu:       &sync.RWMutex{},
		switchMu: &sync.Mutex{},
	}
}
