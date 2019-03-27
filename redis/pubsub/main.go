package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main(){
	go pub()
	go sub()
	sub()
}
func pub(){
	client:=Dial()
	for i:=0;i<1000;i++{
		if err:=client.Publish("test",i).Err();err!=nil{
			panic(err.Error())
		}
	}
}

func sub(){
	client:=Dial()

	sub:=client.Subscribe("test")
	ch:=sub.Channel()
	for msg:=range ch{
		fmt.Printf("收到%s\n",msg.String())
	}
}

func Dial() *redis.Client{
	client:= redis.NewClient(&redis.Options{
		Addr:"localhost:6379",
		DB:1,
	});
	if err:=client.Ping().Err();err!=nil{
		panic(err.Error())
	}
	return client
}