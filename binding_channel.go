package hare

import (
	"github.com/streadway/amqp"
)

type BindingChannel interface {
	BindQueue(*QueueBindingDefinition) error
}

func (c *channel) BindQueue(d *QueueBindingDefinition) error {
	return c.invoke(func(ch *amqp.Channel) error {
		return ch.QueueBind(
			d.Queue.Name(),
			d.Key,
			d.Exchange.Name(),
			d.NoWait,
			d.Args,
		)
	})
}
