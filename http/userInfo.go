package main

import (
	"fmt"
	"net/url"
)

func main() {
	address := "http://user:pass@localhost:1234/debug/pprof"
	u, err := url.Parse(address)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(u.User)
	}
}
