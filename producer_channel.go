package hare

import (
	"github.com/streadway/amqp"
)

type Publishing struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
	Message    amqp.Publishing
}

type ProducerChannel interface {
	Publish(Publishing) error
}

func (c *channel) Publish(p Publishing) error {
	return c.invoke(func(ch *amqp.Channel) error {
		return ch.Publish(p.Exchange, p.RoutingKey, p.Mandatory, p.Immediate, p.Message)
	})
}
