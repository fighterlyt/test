package main

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func main() {
	type Msg struct {
		Usage0 *big.Float
		Usage1 *big.Float
		Usage2 *big.Float
	}

	jsonMsg := `{"Usage0": "31241.4543", "Usage1": "54354325423.65", "Usage2": "123456789012.12"}`

	var msg Msg
	err := json.Unmarshal([]byte(jsonMsg), &msg)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("%+v\n", msg)
	}
}
