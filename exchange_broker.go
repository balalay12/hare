package hare

import "sync"

type ExchangeOption func(*exchange)

type ExchangeBroker interface {
	Exchange(string, ...ExchangeOption) Exchange
}

func (b *broker) Exchange(name string, opts ...ExchangeOption) Exchange {
	return b.exchanges.GetOrSet(name, func() *exchange {
		instance := newExchange(name, opts...)
		instance.channelInvoker = b.Connection().Channel
		return instance
	})
}

func newExchange(name string, opts ...ExchangeOption) *exchange {
	instance := &exchange{
		definition: ExchangeDefinition{
			Name: name,
			Kind: defaultExchangeKind,
		},
		mu: &sync.Mutex{},
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instance)
		}
	}
	return instance
}

func WithExchangeKind(kind string) ExchangeOption {
	return func(e *exchange) {
		e.definition.Kind = kind
	}
}

func WithExchangeDurable() ExchangeOption {
	return func(e *exchange) {
		e.definition.Durable = true
	}
}
