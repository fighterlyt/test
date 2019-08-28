package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
)

var (
	port = "12345"
)

func main() {
	flag.StringVar(&port, "port", port, "监听端口")
	flag.Parse()

	http.HandleFunc("/test", func(resp http.ResponseWriter, r *http.Request) {
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("dump请求错误:" + err.Error())
		} else {
			log.Println("请求:" + string(data))
		}
		resp.Write(data)

	})
	log.Println("启动" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
