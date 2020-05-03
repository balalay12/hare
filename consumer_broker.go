package hare

import "sync"

type ConsumerOption func(*consumer)

type ConsumerBroker interface {
	Consumer(string, ...ConsumerOption) Consumer
}

func (b *broker) Consumer(name string, opts ...ConsumerOption) Consumer {
	return b.consumers.GetOrSet(name, func() *consumer {
		instance := newConsumer(name, opts...)
		instance.channelInvoker = b.Connection().Channel
		return instance
	})
}

func newConsumer(name string, opts ...ConsumerOption) *consumer {
	instance := &consumer{
		definition: ConsumerDefinition{
			Name: name,
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
