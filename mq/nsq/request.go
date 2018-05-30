package main

import (
	"github.com/nats-io/go-nats"
	"fmt"
	"time"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	nc.Subscribe("help", func(m *nats.Msg) {
		m.Data=append(m.Data,'1')
		nc.Publish(m.Reply, m.Data)
	})
	nc.Subscribe("help", func(m *nats.Msg) {
		fmt.Println("heard", string(m.Data))
	})
	// Requests
	msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("response", string(msg.Data))
	}
	// Replies

	// Close connection
	nc, _ = nats.Connect("nats://localhost:4222")
	nc.Close();
}
