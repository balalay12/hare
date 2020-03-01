package main

import (
	"fmt"
	"hare"
)

func main() {
	broker := hare.NewRabbit(
		hare.WithAddr("amqp://guest:guest@127.0.0.1:5672"),
		hare.WithExchange("publisher"),
		hare.WithCount(1),
		hare.WithGlobal(),
	)

	if err := broker.Connect(); err != nil {
		panic(err)
	}

	done := make(chan struct{})

	subOpts := []hare.SubscribeOption{hare.WithQueue("sub")} //, hare.WithDurableQueue()}

	go func() {
		fn := func(msg []byte) error {
			fmt.Println("first", string(msg))

			return nil
		}

		_, err := broker.Subscribe("", fn, subOpts...)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		fn := func(msg []byte) error {
			fmt.Println("second", string(msg))

			return nil
		}

		_, err := broker.Subscribe("", fn, subOpts...)
		if err != nil {
			panic(err)
		}
	}()

	<-done
}
