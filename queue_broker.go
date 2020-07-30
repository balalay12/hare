package hare

import "sync"

type QueueOption func(*queue)

type QueueBroker interface {
	Queue(string, ...QueueOption) Queue
}

func (b *broker) Queue(name string, opts ...QueueOption) Queue {
	return b.queues.GetOrSet(name, func() *queue {
		instance := newQueue(name, opts...)
		instance.channelInvoker = b.Connection().Channel
		return instance
	})
}

func newQueue(name string, opts ...QueueOption) *queue {
	instance := &queue{
		definition: QueueDefinition{
			Name: name,
		},
		mu: &sync.Mutex{},
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instance)
		}
	}
	return instance
}

func WithQueueDurable() QueueOption {
	return func(q *queue) {
		q.definition.Durable = true
	}
}

func WithQueueAutoDelete() QueueOption {
	return func(q *queue) {
		q.definition.AutoDelete = true
	}
}
