package main

import (
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"time"
)

var (
	json  = jsoniter.ConfigCompatibleWithStandardLibrary
	slice = [][]string{
		{
			"a", "b",
		},
		{
			"c", "d",
		},
	}
)

//对于结构的字段,提前序列化和整体序列化相同
func main() {
	data, _ := json.Marshal(slice)

	a := general{
		A: "a",
		B: slice,
	}
	b := before{
		A: "a",
		B: jsoniter.RawMessage(data),
	}

	c := appendStruct{
		A: "a",
	}
	dataA, _ := json.Marshal(a)

	dataB, _ := json.Marshal(b)

	dataC, _ := json.Marshal(c)
	dataC=append(dataC[:len(dataC)-1],[]byte(`,"b":`)...)
	dataC=append(dataC,data...)
	dataC=append(dataC,[]byte(`}`)...)
	if !bytes.Equal(dataA, dataB) || !bytes.Equal(dataA,dataC) {
		fmt.Println("不相同")
		fmt.Println(string(dataA))
		fmt.Println(string(dataB))
		fmt.Println(string(dataC))
	} else {
		timing(10000)
	}
}

type general struct {
	A string     `json:"a"`
	B [][]string `json:"b"`
}

type before struct {
	A string              `json:"a"`
	B jsoniter.RawMessage `json:"b"`
}

func timing(count int, ) {
	start := time.Now()
	var data, dataA, dataB,dataC []byte
	for i := 0; i < count; i++ {
		data, _ = json.Marshal(slice)
	}

	beforeTakes := time.Since(start)

	a := general{
		A: "a",
		B: slice,
	}
	b := before{
		A: "a",
		B: jsoniter.RawMessage(data),
	}
	start = time.Now()
	for i := 0; i < count; i++ {
		dataA, _ = json.Marshal(a)
	}

	allTakes := time.Since(start)

	start = time.Now()
	for i := 0; i < count; i++ {
		dataB, _ = json.Marshal(b)
	}

	seperateTakes := time.Since(start)

	seperateTakes += beforeTakes
	c := appendStruct{
		A: "a",
	}

	start=time.Now()
	for i:=0;i<count;i++{
		dataC, _ = json.Marshal(c)
		dataC=append(dataC[:len(dataC)-1],[]byte(`,"b":`)...)
		dataC=append(dataC,data...)
		dataC=append(dataC,[]byte(`}`)...)
	}
	appendTakes:=time.Since(start)
	appendTakes+=beforeTakes
	fmt.Printf("相同,一起耗时%s,分开耗时%s,附加耗时%s,其中内部耗时%s\n", allTakes.String(), seperateTakes.String(), appendTakes.String(),beforeTakes.String())
	println(string(dataA), string(dataB),string(dataC))
}

type appendStruct struct {
	A string `json:"a"`
}
