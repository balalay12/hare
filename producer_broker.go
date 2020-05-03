package hare

import "sync"

type ProducerOption func(*producer)

type ProducerBroker interface {
	Producer(...ProducerOption) Producer
}

func (b *broker) Producer(opts ...ProducerOption) Producer {
	instance := &producer{
		channelInvoker: b.Connection().Channel,
		channelOptions: make([]ChannelOption, 0),
		mu:             &sync.Mutex{},
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instance)
		}
	}
	return instance
}

func WithProducerChannelOptions(opts ...ChannelOption) ProducerOption {
	return func(p *producer) {
		p.channelOptions = opts
	}
}
