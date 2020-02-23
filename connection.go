package hare

import (
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type connection struct {
	Connection      *amqp.Connection
	Channel         *channel
	ExchangeChannel *channel

	exchange Exchange
	addr     string
	count    int
	global   bool

	sync.Mutex
	connected bool // уже есть соединение
	close     chan struct{}

	waitConn chan struct{}
}

type Exchange struct {
	Name    string
	Durable bool
}

func NewConnection(ex Exchange, addr string, count int, global bool) *connection {
	conn := connection{
		exchange: ex,
		addr:     addr,
		count:    count,
		global:   global,
		close:    make(chan struct{}),
		waitConn: make(chan struct{}),
	}
	close(conn.waitConn)

	return &conn
}

func (c *connection) Connect() (err error) {
	c.Lock()

	if c.connected {
		c.Unlock()
		return
	}

	select {
	case <-c.close:
		c.close = make(chan struct{})
	default:
	}

	c.Unlock()

	return c.connect()
}

func (c *connection) Close() error {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.close:
		return nil
	default:
		close(c.close)
		c.connected = false
	}

	return c.Connection.Close()
}

func (c *connection) Consume(queue, key string, headers, args amqp.Table, autoAck bool) (*channel, <-chan amqp.Delivery, error) {
	consumerCh, err := newChannel(c.Connection, c.count, c.global)
	if err != nil {
		return nil, nil, err
	}

	// TODO durable queue
	if err := consumerCh.DeclareDurableQueue(queue, args); err != nil {
		return nil, nil, err
	}

	deliveries, err := consumerCh.ConsumeQueue(queue, autoAck)
	if err != nil {
		return nil, nil, err
	}

	if err := consumerCh.BindQueue(queue, key, c.exchange.Name, headers); err != nil {
		return nil, nil, err
	}

	return consumerCh, deliveries, nil
}

func (c *connection) Publish(exchange, key string, msg amqp.Publishing) error {
	return c.ExchangeChannel.Publish(exchange, key, msg)
}

func (c *connection) connect() (err error) {
	if err = c.tryConnect(); err != nil {
		return
	}

	c.Lock()
	c.connected = true
	c.Unlock()

	go c.reconnect()

	return
}

func (c *connection) reconnect() {
	var connect bool

	for {
		if connect {
			if err := c.tryConnect(); err != nil {
				time.Sleep(1 * time.Second) // TODO: вынести в конфиг
				continue
			}

			c.Lock()
			c.connected = true
			c.Unlock()

			close(c.waitConn)
		}

		connect = true
		notifyClose := make(chan *amqp.Error)
		c.Connection.NotifyClose(notifyClose)

		select {
		case <-notifyClose:
			c.Lock()
			c.connected = false
			c.waitConn = make(chan struct{})
			c.Unlock()
		case <-c.close:
			return
		}
	}
}

func (c *connection) tryConnect() (err error) {
	// FIXME: подумать над конфигом
	if c.Connection, err = amqp.DialConfig(c.addr, amqp.Config{}); err != nil {
		return
	}

	if c.Channel, err = newChannel(c.Connection, c.count, c.global); err != nil {
		return
	}

	// TODO: durable exchange
	if err = c.Channel.DeclareExchange(c.exchange.Name); err != nil {
		return
	}

	if c.ExchangeChannel, err = newChannel(c.Connection, c.count, c.global); err != nil {
		return
	}

	return
}
