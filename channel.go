package hare

import (
	"github.com/streadway/amqp"
	"sync"
)

type ChannelOption func(*channel)

type ChannelDeclaration interface {
	Declare(...ChannelOption) error
}

type Channel interface {
	ExchangeChannel
	ProducerChannel
	QueueChannel
	ConsumerChannel
	BindingChannel
	IsClosed() bool
	Close() error
	ID() string
}

type channel struct {
	id  string
	qos struct {
		prefetchCount int
		prefetchSize  int
		global        bool
	}

	connected bool

	amqp        *amqp.Channel
	amqpInvoker func() (*amqp.Channel, error)

	mu *sync.Mutex

	notifyClose chan *amqp.Error

	lastError error
}

func (c *channel) ID() string {
	return c.id
}

func (c *channel) Close() error {
	return c.amqp.Close()
}

func (c *channel) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isAmqpClosed()
}

func (c *channel) isAmqpClosed() bool {
	return !c.connected
}

func (c *channel) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isAmqpClosed() {
		return nil
	}
	c.lastError = nil
	go c.checkClose()
	ch, err := c.amqpInvoker()
	if err != nil {
		return err
	}
	ch.NotifyClose(c.notifyClose)
	c.amqp = ch
	if err = c.applyQos(); err != nil {
		return err
	}
	c.connected = true
	return nil
}

func (c *channel) applyQos() error {
	return c.amqp.Qos(c.qos.prefetchCount, c.qos.prefetchSize, c.qos.global)
}

func (c *channel) checkClose() {
	for {
		select {
		case err := <-c.notifyClose:
			c.mu.Lock()
			c.connected = false
			c.amqp = nil
			if err != nil {
				c.lastError = err
			}
			c.mu.Unlock()
			return
		}
	}
}

func (c *channel) invoke(handler func(*amqp.Channel) error) error {
	c.mu.Lock()
	if c.lastError != nil {
		c.mu.Unlock()
		return c.lastError
	}
	c.mu.Unlock()
	if c.IsClosed() {
		if err := c.connect(); err != nil {
			return err
		}
	}
	return handler(c.amqp)
}
