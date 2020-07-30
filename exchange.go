package hare

import (
	"github.com/streadway/amqp"
	"sync"
)

const defaultExchangeKind = "direct"

type ExchangeDefinition struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type Exchange interface {
	ChannelDeclaration
	Name() string
	IsDeclared() bool
}

type exchange struct {
	definition     ExchangeDefinition
	channelInvoker ChannelInvoker

	declared bool
	mu       *sync.Mutex
}

func (e *exchange) Name() string {
	return e.definition.Name
}

func (e *exchange) IsDeclared() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.declared
}
