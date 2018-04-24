package main

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func main() {
	b := big.NewInt(11111111110)
	text, err := b.MarshalText()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(string(text))
	}

	if err = json.Unmarshal([]byte("123456"), b); err != nil {
		panic(err.Error())
	} else {
		text, err = b.MarshalText()
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println(string(text))
		}
	}
}
