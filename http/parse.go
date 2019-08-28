package main

import (
	"code.aliyun.com/soyou/platform/tools"
	"github.com/davecgh/go-spew/spew"
	"net/url"
	"strings"
)

var data = `table=dealDay%3D20190819%40%23%40dealTime%3D16%3A32%3A10%40%23%40time%3D20190819163210%40%23%40merchantId%3D898210252111883%40%23%40dealMoney%3D0.10%40%23%40discountMoney%3D0.00%40%23%40indexNum%3D09655823930N%40%23%40payType%3D%E5%BE%AE%E4%BF%A1%40%23%40serialNum%3D%40%23%40dealType%3D%E6%88%90%E5%8A%9F%40%23%40dealStatus%3D%E6%88%90%E5%8A%9F%40%23%40terminal%3D19071521%40%23%40%7C-%7C`

func main() {
	values, err := url.ParseQuery(data)
	if err != nil {
		panic(err.Error())
	}
	spew.Dump(values)
	for key, keyValues := range values {
		for _, value := range keyValues {

			fields := strings.Split(value, "@#@")
			for _, field := range fields {
				eachFields := strings.Split(field, "=")
				if len(eachFields) == 2 {

				}
			}
		}
	}
}

type ListImportResult struct {
	List []*ImportResult `form:"list"`
}

type ImportResult struct {
	Abstract    string `form:"abstract"`
	AccountName string `form:"accountName"`
	AccountNum  string `form:"accountNum"`
	AccountText string `form:"accountText"`
	Balance     string `form:"balance"`
	In          string `form:"in"`
	Out         string `form:"out"`
	Time        string `form:"time"`
	Unit        string `form:"unit"`
	Order       string `form:"order"`
}

type Deal struct {
	DealTime   tools.Time
	Money      string
	Number     string //商户编号
	Transplace string //支付方式
	DealType   string //交易类型

	DealStatus string //交易状态
}
