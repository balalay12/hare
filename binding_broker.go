package hare

import "github.com/streadway/amqp"

type BindingBroker interface {
	BindQueue(Queue, Exchange, ...QueueBindingOption) error
}

type QueueBindingOption func(*QueueBindingDefinition)

type QueueBindingDefinition struct {
	Queue
	Exchange
	Key    string
	NoWait bool
	Args   amqp.Table
}

func (b *broker) BindQueue(queue Queue, exchange Exchange, opts ...QueueBindingOption) error {
	definition := &QueueBindingDefinition{
		Queue:    queue,
		Exchange: exchange,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(definition)
		}
	}
	return b.Connection().Channel().BindQueue(definition)
}

func WithQueueBindingKey(key string) QueueBindingOption {
	return func(definition *QueueBindingDefinition) {
		definition.Key = key
	}
}

func WithQueueBindingNoWait() QueueBindingOption {
	return func(definition *QueueBindingDefinition) {
		definition.NoWait = true
	}
}

func WithQueueBindingArgs(args amqp.Table) QueueBindingOption {
	return func(definition *QueueBindingDefinition) {
		definition.Args = args
	}
}
