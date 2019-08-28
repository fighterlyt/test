package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	for i := 0; i < 1001; i++ {
		send(i)
	}
	time.Sleep(time.Second)
}
func send(i int) {
	// 查询卡状态 (POST http://localhost:1234/api/invoke)

	body := strings.NewReader(`{ "id":"5d09dfa65bf6c8bb15cfed5f",
	"time":1561383052,
	"key":"123",
	"apiKey":"card.status",
	"data":["6217003320075998396"]
}`)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://localhost:1234/api/invoke", body)

	// Headers
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	if resp.Body != nil {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body : ", string(respBody), i)
	}

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)

}
