package hare

import (
	"github.com/streadway/amqp"
)

type ExchangeChannel interface {
	DeclareExchange(*ExchangeDefinition) error
}

func (e *exchange) Declare(opts ...ChannelOption) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.declared {
		return nil
	}
	ch := e.channelInvoker(opts...)
	defer func() {
		_ = ch.Close()
	}()
	err := ch.DeclareExchange(&e.definition)
	if err != nil {
		return err
	}
	e.declared = true
	return nil
}

func (c *channel) DeclareExchange(d *ExchangeDefinition) error {
	return c.invoke(func(ch *amqp.Channel) error {
		return ch.ExchangeDeclare(
			d.Name,
			d.Kind,
			d.Durable,
			d.AutoDelete,
			d.Internal,
			d.NoWait,
			d.Args,
		)
	})
}
