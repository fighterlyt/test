package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"io/ioutil"
)

var (
	path = ""
)

func main() {
	flag.StringVar(&path, "path", "", "文件路径")
	flag.Parse()

	if path==""{
		path="/Users/mac/go/src/github.com/fighterlyt/test/compress/gzip/data/quote.json"
	}
	if data, err := ioutil.ReadFile(path); err != nil {
		panic(err.Error())
	} else {

		buffer := &bytes.Buffer{}
		writer, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
		if err != nil {
			panic("构建 gzip.Writer 失败" + err.Error())
		}

		if _, err := writer.Write(data); err != nil {
			panic(err.Error())
		}

		writer.Flush()
		writer.Close()
		println(len(data), buffer.Len())

		gReader,err:=gzip.NewReader(buffer)
		if err!=nil{
			panic(err.Error())
		}
		readData,err:=ioutil.ReadAll(gReader)
		if err!=nil{
			panic(err.Error())
		}
		println(string(readData))
	}
}
