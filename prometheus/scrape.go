package main

import (
	"net/http"
	"fmt"
	"io"
	"bytes"
	"compress/gzip"
)

const acceptHeader = `text/plain;version=0.0.4;q=1,*/*;q=0.1`

var userAgentHeader = fmt.Sprintf("Prometheus/%s", "1.0")

func main() {
	w := &bytes.Buffer{}
	req, err := http.NewRequest("GET", "http://localhost:6060/metrics", nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Add("Accept", acceptHeader)
	//req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgentHeader)
	req.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", fmt.Sprintf("%f", 5.5))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("headers",resp.Header)
	if resp.StatusCode != http.StatusOK {
		panic("server returned HTTP status %s" + resp.Status)
	}

	if resp.Header.Get("Content-Encoding") != "gzip" {
		_, err = io.Copy(w, resp.Body)
		if err!=nil{
			panic(err.Error())
		}else{
			fmt.Println(w.String())
		}
	}else{
		gzipr, err := gzip.NewReader(resp.Body)
		if err != nil {
			panic(err.Error())
		}
		//if data,err:=httputil.DumpResponse(resp,true);err!=nil{
		//	panic(err.Error())
		//}else{
		//fmt.Println("resp",string(data))
		//
		//}

		_, err = io.Copy(w, gzipr)
		gzipr.Close()
		if err!=nil{
			panic(err.Error())
		}else{
			fmt.Println(w.String())
		}
	}


}
