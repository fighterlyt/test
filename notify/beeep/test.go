package main

import "github.com/gen2brain/beeep"

func main() {
	err := beeep.Alert("警告", "Message body", "assets/information.png")
	if err != nil {
		panic(err)
	}
}
