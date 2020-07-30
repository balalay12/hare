package hare

import (
	"github.com/streadway/amqp"
	"sync"
)

type Connection interface {
	ChannelBroker
	IsClosed() bool
}

type connection struct {
	dsn string

	amqp     *amqp.Connection
	channels map[string]*channel

	mu     *sync.Mutex
	chanMu *sync.RWMutex

	notifyClose chan *amqp.Error
}

func (c *connection) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isAmqpClosed()
}

func (c *connection) isAmqpClosed() bool {
	return c.amqp == nil || c.amqp.IsClosed()
}

func (c *connection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isAmqpClosed() {
		return nil
	}
	go c.checkState()
	internal, err := amqp.Dial(c.dsn)
	if err != nil {
		return err
	}
	internal.NotifyClose(c.notifyClose)
	c.amqp = internal
	return nil
}

func (c *connection) checkState() {
	for {
		select {
		case <-c.notifyClose:
			c.amqp = nil
			return
		}
	}
}

func (c *connection) invoke(handler func(*amqp.Connection) error) error {
	if c.IsClosed() {
		if err := c.connect(); err != nil {
			return err
		}
	}
	return handler(c.amqp)
}
