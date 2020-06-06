package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// cURL (POST https://59bfc73f0956497c8fe3e1e3b8ca9d2c.ap-northeast-1.aws.found.io:9243/br-cost_orders_logs/_search)

	json := []byte(`{"query": {"range": {"add_time": {"gt": "now-769h","lte": "now-768h"}}},"size": 1}`)
	body := bytes.NewBuffer(json)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://59bfc73f0956497c8fe3e1e3b8ca9d2c.ap-northeast-1.aws.found.io:9243/br-cost_orders_logs/_search", body)

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic ZWxhc3RpYzpVYUZHbnhzdEpjeDhtWEZ2UlJSa050ZjU=")
	req.Header.Add("Accept-Encoding", "gzip")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	reader, _ := gzip.NewReader(resp.Body)
	// Read Response Body
	respBody, _ := ioutil.ReadAll(reader)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}
