package hare

import (
	"errors"
	"github.com/streadway/amqp"
)

type ConsumerChannel interface {
	Consume(Queue, *ConsumerDefinition, []ConsumerHandler) error
}

func (c *consumer) Consume(queue Queue) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.handlers) == 0 {
		return errors.New("no handlers were defined yet")
	}
	return c.invoke(func(ch ConsumerChannel) error {
		if !queue.IsDeclared() {
			err := queue.Declare()
			if err != nil {
				return err
			}
		}
		return ch.Consume(queue, &c.definition, c.handlers)
	})
}

func (c *channel) Consume(queue Queue, definition *ConsumerDefinition, handlers []ConsumerHandler) error {
	return c.invoke(func(ch *amqp.Channel) error {
		messages, err := ch.Consume(
			queue.Name(),
			definition.Name,
			definition.AutoAck,
			definition.Exclusive,
			definition.NoLocal,
			definition.NoWait,
			definition.Args,
		)
		if err != nil {
			return err
		}
		for message := range messages {
			msg := &consumerMessage{
				amqp: &message,
			}
			for _, handler := range handlers {
				handler(msg)
			}
			if msg.needReQueue() {
				err := message.Nack(true, true)
				if err != nil {
					return err
				}
			} else {
				if msg.isAckQuorum(0.0) {
					err := message.Ack(true)
					if err != nil {
						return err
					}
				} else {
					err := message.Nack(true, false)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}
