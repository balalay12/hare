package hare

import (
	"errors"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var ErrNilChannel = errors.New("channel is nil")

type channel struct {
	uuid    string
	conn    *amqp.Connection
	channel *amqp.Channel
}

func newChannel(conn *amqp.Connection, count int, global bool) (*channel, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	channel := channel{
		uuid: id.String(),
		conn: conn,
	}

	if err := channel.Connect(count, global); err != nil {
		return nil, err
	}

	return &channel, nil
}

// FIXME make it unexported
func (c *channel) Connect(count int, global bool) (err error) {
	if c.channel, err = c.conn.Channel(); err != nil {
		return
	}

	if err = c.channel.Qos(count, 0, global); err != nil {
		return
	}

	return
}

// FIXME make it unexported
func (c *channel) Close() error {
	if c.channel == nil {
		return ErrNilChannel
	}

	return c.channel.Close()
}

func (c *channel) Publish(ex, key string, msg amqp.Publishing) error {
	if c.channel == nil {
		return ErrNilChannel
	}

	return c.channel.Publish(ex, key, false, false, msg)
}

func (c *channel) DeclareExchange(exchange string) error {
	return c.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // kind
		false,    // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
}

func (c *channel) DeclareDurableExchange(exchange string) error {
	return c.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
}

func (c *channel) DeclareQueue(queue string, args amqp.Table) error {
	_, err := c.channel.QueueDeclare(
		queue, // name
		false, // durable
		true,  // autoDelete
		false, // exclusive
		false, // noWait
		args,  // args
	)
	return err
}

func (c *channel) DeclareDurableQueue(queue string, args amqp.Table) error {
	_, err := c.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		args,  // args
	)
	return err
}

func (c *channel) DeclareReplyQueue(queue string) error {
	_, err := c.channel.QueueDeclare(
		queue, // name
		false, // durable
		true,  // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

func (c *channel) ConsumeQueue(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue,   // queue
		c.uuid,  // consumer
		autoAck, // autoAck
		false,   // exclusive
		false,   // nolocal
		false,   // nowait
		nil,     // args
	)
}

func (c *channel) BindQueue(queue, key, exchange string, args amqp.Table) error {
	return c.channel.QueueBind(
		queue,    // name
		key,      // key
		exchange, // exchange
		false,    // noWait
		args,     // args
	)
}
