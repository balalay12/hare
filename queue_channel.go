package hare

import (
	"github.com/streadway/amqp"
)

type QueueChannel interface {
	DeclareQueue(*QueueDefinition, *amqp.Queue) error
}

func (q *queue) Declare(opts ...ChannelOption) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.declared {
		return nil
	}
	ch := q.channelInvoker(opts...)
	defer func() {
		_ = ch.Close()
	}()
	err := ch.DeclareQueue(&q.definition, q.amqp)
	if err != nil {
		return err
	}
	q.declared = true
	return nil
}

func (c *channel) DeclareQueue(d *QueueDefinition, queue *amqp.Queue) error {
	return c.invoke(func(ch *amqp.Channel) error {
		q, err := ch.QueueDeclare(
			d.Name,
			d.Durable,
			d.AutoDelete,
			d.Exclusive,
			d.NoWait,
			d.Args,
		)
		if err != nil {
			return err
		}
		queue = &q
		return nil
	})
}
