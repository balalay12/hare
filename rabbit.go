package hare

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var ErrNoConnection = errors.New("connection is nil")

type rabbit struct {
	conn *connection
	addr string
	opts Options
	wg   sync.WaitGroup
	mtx  sync.Mutex
}

type subscriber struct {
	mtx          sync.Mutex
	mayRun       bool
	opts         SubscribeOptions
	routingKey   string
	durableQueue bool
	args         map[string]interface{}
	r            *rabbit
	ch           *channel
	callback     func(msg amqp.Delivery)
}

func NewRabbit(opts ...Option) *rabbit {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return &rabbit{
		addr: opt.addr,
		opts: opt,
	}
}

func (r *rabbit) Publish(topic string, msg []byte) error {
	/*
		- Delivery mode
		- Header
	*/

	m := amqp.Publishing{
		Body: msg,
	}

	if r.conn == nil {
		return errors.New("connection is nil")
	}

	return r.conn.Publish(r.conn.exchange.Name, topic, m)
}

func (r *rabbit) Subscribe(topic string, handler func(msg []byte) error, opts ...SubscribeOption) (*subscriber, error) {
	//var ackSuccess bool

	if r.conn == nil {
		return nil, errors.New("no connection")
	}

	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	fn := func(m amqp.Delivery) {
		err := handler(m.Body)
		if err != nil {
			m.Ack(false)
		} else {
			m.Nack(false, false) // TODO: requeue opts
		}
	}

	s := subscriber{
		opts:       opt,
		routingKey: topic,
		r:          r,
		callback:   fn,
		mayRun:     true,
	}

	go s.resubscribe()

	return &s, nil
}

func (s *subscriber) Unsubscribe() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.mayRun = false

	if s.ch != nil {
		return s.ch.Close()
	}

	return nil
}

func (r *rabbit) Connect() error {
	ex := Exchange{Name: r.opts.exchange} // FIXME: durable?

	if r.conn == nil {
		r.conn = NewConnection(ex, r.opts.addr, r.opts.count, r.opts.global)
	}

	return r.conn.connect()
}

func (r *rabbit) Disconnect() error {
	if r.conn == nil {
		return ErrNoConnection
	}

	err := r.conn.Close()
	r.wg.Wait()
	return err
}

func (s *subscriber) resubscribe() {
	minResubscribeDelay := 100 * time.Millisecond
	maxResubscribeDelay := 30 * time.Second
	expFactor := time.Duration(2)
	reSubscribeDelay := minResubscribeDelay

	for {
		s.mtx.Lock()
		mayRun := s.mayRun
		s.mtx.Unlock()
		if !mayRun {
			return
		}

		select {
		case <-s.r.conn.close:
			//
			return
		//wait until we reconect to rabbit
		case <-s.r.conn.waitConn:
		}

		// it may crash (panic) in case of Consume without connection, so recheck it
		s.r.mtx.Lock()
		if !s.r.conn.connected {
			s.r.mtx.Unlock()
			continue
		}

		/*
			FIXME: if create durable queue and forgot option for that we never handle error about it
			Example: Exception (406) Reason: "PRECONDITION_FAILED - inequivalent arg 'durable' for queue 'sub'
				in vhost '/': received 'false' but current is 'true'"
		*/
		ch, sub, err := s.r.conn.Consume(
			s.opts.Queue,
			s.routingKey,
			amqp.Table{}, // FIXME
			amqp.Table{}, // FIXME
			s.opts.AutoAck,
			s.opts.Durable,
		)

		s.r.mtx.Unlock()
		switch err {
		case nil:
			reSubscribeDelay = minResubscribeDelay
			s.mtx.Lock()
			s.ch = ch
			s.mtx.Unlock()
		default:
			if reSubscribeDelay > maxResubscribeDelay {
				reSubscribeDelay = maxResubscribeDelay
			}
			time.Sleep(reSubscribeDelay)
			reSubscribeDelay *= expFactor
			continue
		}
		for d := range sub {
			s.r.wg.Add(1)
			s.callback(d)
			s.r.wg.Done()
		}
	}

}
