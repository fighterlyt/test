package main

import (
	"github.com/songtianyi/wechat-go/wxweb"
	"time"
	"log"
)

func main() {
	session, err := wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	if err != nil {
		panic(err.Error())
		return
	}
	if err := session.LoginAndServe(false); err != nil {
		panic(err.Error())
	}

	log.Println("发送")
	if _, _, err := session.SendText("test", "1", "2"); err != nil {
		panic(err.Error())
	}

	time.Sleep(time.Minute)

}
