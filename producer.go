package hare

import "sync"

type Producer interface {
	ProducerPublishing
}

type producer struct {
	channelInvoker ChannelInvoker
	channelOptions []ChannelOption
	channel        ProducerChannel

	mu *sync.Mutex
}

func (p *producer) invoke(handler func(c ProducerChannel) error) error {
	if p.channel == nil {
		p.channel = p.channelInvoker(p.channelOptions...)
	}
	return handler(p.channel)
}
