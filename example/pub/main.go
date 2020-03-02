package main

import (
	"strconv"
	"time"

	"github.com/balalay12/hare"
)

func main() {
	broker := hare.NewRabbit(
		hare.WithAddr("amqp://guest:guest@127.0.0.1:5672"),
		hare.WithExchange("publisher"),
		hare.WithDurableExchange(),
		hare.WithCount(1),
	)

	if err := broker.Connect(); err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		if err := broker.Publish("", []byte("test "+strconv.Itoa(i))); err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
