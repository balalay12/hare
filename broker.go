package hare

type BrokerOption func(*broker)

type Broker interface {
	ConnectionBroker
	ChannelBroker
	ExchangeBroker
	ProducerBroker
	QueueBroker
	ConsumerBroker
	BindingBroker
}

type broker struct {
	autoConnect bool

	conn      *connection
	exchanges *exchangeCollection
	queues    *queueCollection
	consumers *consumerCollection
}

func NewBroker(dsn string, opts ...BrokerOption) (Broker, error) {
	instance := &broker{
		exchanges: newExchangeCollection(),
		queues:    newQueueCollection(),
		consumers: newConsumerCollection(),
		conn:      newConnection(dsn),
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instance)
		}
	}
	if instance.autoConnect {
		err := instance.conn.connect()
		if err != nil {
			return instance, err
		}
	}
	return instance, nil
}

func WithAutoConnect() BrokerOption {
	return func(b *broker) {
		b.autoConnect = true
	}
}
