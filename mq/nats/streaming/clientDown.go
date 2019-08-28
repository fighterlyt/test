package main

import (
	"github.com/fighterlyt/test/mq/nats/streaming/base"
	stan "github.com/nats-io/go-nats-streaming"
	"strconv"
	"sync"
	"time"
)

func main() {

	finish := &sync.WaitGroup{}
	ready := &sync.WaitGroup{}
	finish.Add(2)
	ready.Add(1)

	go client(finish, ready)
	ready.Wait()

	go server(finish)
	finish.Wait()
}

func server(finish *sync.WaitGroup) {
	conn, err := base.Connect("nats://localhost:4222", "test1")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	println("生产者启动")
	for i := 0; i < 10; i++ {
		if err := conn.Publish("test", []byte(strconv.Itoa(i))); err != nil {
			panic(err.Error())
		}

	}
	finish.Done()
	println("生产者结束")

}

func client(finish, ready *sync.WaitGroup) {
	conn, err := base.Connect("nats://localhost:4222", "test")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	retry := make(chan struct{}, 1)
	var sub stan.Subscription
	var firstProcess func(msg *stan.Msg)
	firstProcess = func(msg *stan.Msg) {
		msg.Ack()
		if string(msg.Data) == "3" {
			println("消费者停止订阅,1S后重新订阅")
			sub.Close()

			time.Sleep(time.Second)

			retry <- struct{}{}
		}

		println(string(msg.Data))

	}
	sub, err = conn.Subscribe("test", firstProcess)

	if err != nil {
		panic("订阅失败" + err.Error())
	}
	ready.Done()
	println("消费者完成订阅")

	<-retry
	sub, err = conn.Subscribe("test", func(msg *stan.Msg) {
		println(string(msg.Data), string(msg.Data) == "9")
		msg.Ack()
		if string(msg.Data) == "9" {
			finish.Done()
			println("消费者结束")
		}

	})

	if err != nil {
		panic("订阅失败" + err.Error())
	}
	println("消费者重新订阅完成")

}
