package main

import (
	"os"
	"time"
)

const (
	key = "BIDHOUSE"
)

func main() {
	if _, exist := os.LookupEnv(key); !exist {
		println("设置,启动", os.Getenv(key))

		os.Setenv(key, "running")
		println("设置,启动", os.Getenv(key))

	} else {
		println("等待")
		for {
			time.Sleep(time.Second)
			if os.Getenv(key) == "" {
				println("设置,启动", os.Getenv(key))

				os.Setenv(key, "running")
				println("设置,启动", os.Getenv(key))

			}
			println("等待")
		}
	}
	defer os.Unsetenv(key)
	time.Sleep(time.Second * 10)
	println("退出")
}
