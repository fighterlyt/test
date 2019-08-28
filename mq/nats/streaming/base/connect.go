package base

import (
	stan "github.com/nats-io/go-nats-streaming"
)

func Connect(url, name string) (stan.Conn, error) {
	if conn, err := stan.Connect("test-cluster", name, stan.NatsURL(url)); err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}
