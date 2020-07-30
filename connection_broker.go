package hare

import (
	"github.com/streadway/amqp"
	"sync"
)

type ConnectionBroker interface {
	Connection() Connection
}

func (b *broker) Connection() Connection {
	return b.conn
}

func newConnection(dsn string) *connection {
	return &connection{
		dsn:         dsn,
		channels:    make(map[string]*channel),
		mu:          &sync.Mutex{},
		chanMu:      &sync.RWMutex{},
		notifyClose: make(chan *amqp.Error),
	}
}
