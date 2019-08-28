package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"radar-trade/help"
	"time"
)

const (
	host     = "http://47.52.254.119:19301/"
	sendurl  = "sms/toSent"
	poolsend = "__catpool__/catpool/pull"
	key      = "abc123"
)

func main() {
	err := send("18665597023", "18665735386", "134")
	if err != nil {
		panic(err.Error())
	}
}
func send(from, to, content string) error {
	t := time.Now().Unix()
	sign := help.Md5([]byte(fmt.Sprintf("%d%s", t, key)), false)

	argu := SendArgument{
		Time:    t,
		Sign:    sign,
		From:    from,
		Target:  to,
		Content: content,
	}
	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(argu); err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, host+sendurl, buffer)
	if err != nil {
		panic(err.Error())
	}
	reqData, _ := httputil.DumpRequest(req, true)
	println("请求", string(reqData))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		println(string(data))
		return nil
	}

}

type SendArgument struct {
	Time    int64  `json:"time"`
	Sign    string `json:"sign"`
	From    string `json:"from_phone"`
	Target  string `json:"accept_phone"`
	Content string `json:"content"`
}
