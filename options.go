package hare

type Options struct {
	addr     string
	exchange string
	count    int
	global   bool

	Subscriber SubscribeOptions
}

type Option func(options *Options)

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

func WithExchange(ex string) Option {
	return func(o *Options) {
		o.exchange = ex
	}
}

func WithCount(c int) Option {
	return func(o *Options) {
		o.count = c
	}
}

func WithGlobal() Option {
	return func(o *Options) {
		o.global = true
	}
}

type PublishOptions struct{}

type PublishOption func(*PublishOptions)

type SubscribeOptions struct {
	AutoAck bool
	Queue   string
}

type SubscribeOption func(*SubscribeOptions)

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithQueue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}
