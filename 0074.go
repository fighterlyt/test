package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	start := int64(65318)
	file, _ := os.Create("0074.log")
	defer file.Close()
	for i := 0; i < 100; i++ {
		go func() {
			for {
				link := fmt.Sprintf("http://dev.10086.cn/dqpt/cmpp5n/api/getAccessUrl?sspId=%d&pageId=P170808008", atomic.AddInt64(&start, 1))
				fmt.Fprintln(file, "测试", link)
				resp, err := http.Get(link)
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode == 302 || resp.StatusCode == 200 {
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
