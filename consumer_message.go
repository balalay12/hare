package hare

import (
	"github.com/streadway/amqp"
	"math"
)

type ConsumerMessage interface {
	Body() []byte

	ReQueue()
	Ack()
	Nack()
}

type consumerMessage struct {
	amqp *amqp.Delivery

	acks    float64
	nacks   float64
	requeue bool
}

func (c *consumerMessage) needReQueue() bool {
	return c.requeue
}

func (c *consumerMessage) isAckQuorum(threshold float64) bool {
	if c.nacks == 0 && c.acks >= 0 {
		return true
	}
	if threshold > 1.0 {
		threshold = 1.0
	}
	return threshold < (math.Round((c.acks/c.nacks)*100.0) / 100.0)
}

func (c *consumerMessage) Body() []byte {
	return c.amqp.Body
}

func (c *consumerMessage) ReQueue() {
	c.requeue = true
}

func (c *consumerMessage) Ack() {
	c.acks++
}

func (c *consumerMessage) Nack() {
	c.nacks++
}
