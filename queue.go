package hare

import (
	"github.com/streadway/amqp"
	"sync"
)

type QueueDefinition struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type Queue interface {
	ChannelDeclaration
	Name() string
	IsDeclared() bool
}

type queue struct {
	definition     QueueDefinition
	channelInvoker ChannelInvoker

	declared bool
	mu       *sync.Mutex

	amqp *amqp.Queue
}

func (q *queue) Name() string {
	return q.definition.Name
}

func (q *queue) IsDeclared() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.declared
}
