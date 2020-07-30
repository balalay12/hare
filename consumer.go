package hare

import (
	"github.com/streadway/amqp"
	"sync"
)

type ConsumerDefinition struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type Consumer interface {
	Name() string
	Handle(...ConsumerHandler) Consumer
	Consume(Queue) error
}

type ConsumerHandler func(ConsumerMessage)

type consumer struct {
	definition     ConsumerDefinition
	channelInvoker ChannelInvoker
	channel        ConsumerChannel

	mu *sync.Mutex

	handlers []ConsumerHandler
}

func (c *consumer) Name() string {
	return c.definition.Name
}

func (c *consumer) Handle(handlers ...ConsumerHandler) Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers = append(c.handlers, handlers...)
	return c
}

func (c *consumer) invoke(handler func(ch ConsumerChannel) error) error {
	if c.channel == nil {
		c.channel = c.channelInvoker()
	}
	return handler(c.channel)
}
