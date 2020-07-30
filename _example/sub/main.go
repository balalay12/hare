package main

import (
	"fmt"
	"github.com/balalay12/hare"
)

func main() {
	broker, err := hare.NewBroker(
		"amqp://guest:guest@127.0.0.1:5672/",
		hare.WithAutoConnect(),
	)
	if err != nil {
		panic(err)
	}
	queue := broker.Queue("subscription")
	consumer := broker.Consumer("subscribe")
	err = consumer.
		Handle(func(message hare.ConsumerMessage) {
			fmt.Println(string(message.Body()))
		}).
		Consume(queue)
	if err != nil {
		panic(err)
	}
}
