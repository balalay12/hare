package hare

import (
	"github.com/streadway/amqp"
)

type PublishingOption func(*publishingOptions)

type publishingOptions struct {
	routingKey string
	headers    amqp.Table
}

type ProducerPublishing interface {
	Publish(Exchange, []byte, ...PublishingOption) error
}

func (p *producer) Publish(exchange Exchange, msg []byte, opts ...PublishingOption) error {
	options := publishingOptions{}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(&options)
		}
	}
	return p.invoke(func(c ProducerChannel) error {
		if !exchange.IsDeclared() {
			err := exchange.Declare()
			if err != nil {
				return err
			}
		}
		return c.Publish(Publishing{
			Exchange:   exchange.Name(),
			RoutingKey: options.routingKey,
			Message: amqp.Publishing{
				Body:    msg,
				Headers: options.headers,
			},
		})
	})
}

func WithPublishRoutingKey(value string) PublishingOption {
	return func(o *publishingOptions) {
		o.routingKey = value
	}
}

func WithPublishHeader(key string, value interface{}) PublishingOption {
	return func(o *publishingOptions) {
		o.headers[key] = value
	}
}
