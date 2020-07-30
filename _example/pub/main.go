package main

import (
	"fmt"
	"github.com/balalay12/hare"
	"strconv"
	"time"
)

func main() {
	broker, err := hare.NewBroker(
		"amqp://guest:guest@127.0.0.1:5672/",
		hare.WithAutoConnect(),
	)
	if err != nil {
		panic(err)
	}
	exchange := broker.Exchange("publisher")
	producer := broker.Producer()
	for i := 0; i < 100; i++ {
		if err := producer.Publish(exchange, []byte("test "+strconv.Itoa(i))); err != nil {
			fmt.Println(fmt.Sprintf("%+v", err))
		} else {
			fmt.Println("Successfully publish")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
