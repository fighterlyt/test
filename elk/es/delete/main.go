package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	count = flag.Int("count", 10, "并发")
	fetch = flag.Int("fetch", 10, "获取数量")
	hour  = flag.Int("hour", 30*24+1, "开始小时")
)

func main() {
	flag.Parse()
	i := new(int64)

	ch := make(chan string, *count*2)
	for j := 0; j < *count; j++ {
		go func() {
			client := &http.Client{}

			for id := range ch {
				// Create client
				req, err := http.NewRequest(http.MethodDelete, "https://59bfc73f0956497c8fe3e1e3b8ca9d2c.ap-northeast-1.aws.found.io:9243/br-cost_orders_logs/br-cost_orders_logs/"+id, nil)
				resp, err := client.Do(req)

				if err != nil {
					return
				}
				count := atomic.AddInt64(i, 1)
				if count%1000 == 0 {
					println(time.Now().Format("2006-01-02 15:04:05")+"已清理", count)
				}
				resp.Body.Close()
			}
		}()
	}

	client := &http.Client{}
	for {
		sendCurl(ch, *fetch, client, 1, *hour)

	}

}
func sendCurl(ch chan string, fetch int, client *http.Client, id, hour int) {
	// cURL (POST https://59bfc73f0956497c8fe3e1e3b8ca9d2c.ap-northeast-1.aws.found.io:9243/br-cost_orders_logs/_search)

	fmt.Printf("[%d]开始获取\n", id)
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"query": {"range": {"add_time": {"gt": "now-%dh","lte": "now-%dh"}}},"size": %d}`, hour+24, hour, fetch)))

	// Create client

	// Create request
	req, err := http.NewRequest("POST", "https://59bfc73f0956497c8fe3e1e3b8ca9d2c.ap-northeast-1.aws.found.io:9243/br-cost_orders_logs/_search", body)

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic ZWxhc3RpYzpVYUZHbnhzdEpjeDhtWEZ2UlJSa050ZjU=")
	req.Header.Add("Accept-Encoding", "gzip")
	// req.Header.Add("Accept-Encoding", "gzip")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	result := &Result{}
	// Display Results
	fmt.Printf("[%d] response Status :%s ,response Headers: %s\n", id, resp.Status, resp.Header)
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	if err := json.NewDecoder(reader).Decode(result); err != nil {
		panic(err.Error())
	}
	if len(result.Hits.Hits) == 0 {
		return
	}
	for _, doc := range result.Hits.Hits {
		ch <- doc.ID
	}

}

type Result struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`

	Hits struct {
		Total    int     `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				Type       int    `json:"type"`
				Message    string `json:"message"`
				ExtMsg     string `json:"ext_msg"`
				AddTime    int    `json:"add_time"`
				ExtOrderID []int  `json:"ext_order_id"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
