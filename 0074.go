package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	start := int64(0)
	file, _ := os.Create("0074.log")
	defer file.Close()

	client := http.Client{}
	client.CheckRedirect = nil
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				link := fmt.Sprintf("http://dev.10086.cn/dqpt/cmpp5n/api/getAccessUrl?sspId=%d&pageId=P170808008", atomic.AddInt64(&start, 1))
				resp, err := client.Get(link)
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode == 302 {
						fmt.Fprintln(file, "可用", link)
					}
				} else {
					// println("不可用", link)
				}
			}

		}()
	}
	time.Sleep(time.Hour)
}
