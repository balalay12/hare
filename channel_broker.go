package hare

import (
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"sync"
)

type ChannelInvoker func(...ChannelOption) Channel

type ChannelBroker interface {
	Channel(...ChannelOption) Channel
}

func (b *broker) Channel(opts ...ChannelOption) Channel {
	return b.conn.Channel(opts...)
}

func (c *connection) Channel(opts ...ChannelOption) Channel {
	c.chanMu.Lock()
	defer c.chanMu.Unlock()
	ch := newChannel(opts...)
	ch.amqpInvoker = c.newAmqpChannel
	c.channels[ch.id] = ch
	return ch
}

func (c *connection) newAmqpChannel() (ch *amqp.Channel, err error) {
	err = c.invoke(func(amqpConn *amqp.Connection) error {
		ch, err = amqpConn.Channel()
		if err != nil {
			ch = nil
			return err
		}
		return nil
	})
	return
}

func newChannel(opts ...ChannelOption) *channel {
	instance := &channel{
		id:          generateChannelId(),
		mu:          &sync.Mutex{},
		notifyClose: make(chan *amqp.Error, 1),
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instance)
		}
	}
	return instance
}

func generateChannelId() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func WithChannelQosPrefetchCount(value int) ChannelOption {
	return func(c *channel) {
		c.qos.prefetchCount = value
	}
}

func WithChannelQosPrefetchSize(value int) ChannelOption {
	return func(c *channel) {
		c.qos.prefetchSize = value
	}
}

func WithChannelQosGlobal() ChannelOption {
	return func(c *channel) {
		c.qos.global = true
	}
}
