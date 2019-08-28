package main

import (
	"bufio"
	"bytes"
	"github.com/dustin/go-humanize"
	"github.com/tdewolff/parse/strconv"
	"io"
	"os"
)

var (
	file = "/Users/mac/memory.log"
)

func main() {
	if file, err := os.Open(file); err != nil {
		panic(err.Error())
	} else {
		defer file.Close()

		reader := bufio.NewReader(file)
		ch := make(chan []byte, 1)

		go func() {
			for {
				data, _, err := reader.ReadLine()
				if err != nil {
					close(ch)
					if err != io.EOF {
						panic(err.Error())
					}
					break

				} else {
					ch <- data
				}
			}

		}()
		rb, tb := process(ch)
		println(humanize.Bytes(uint64(rb)), humanize.Bytes(uint64(tb)))
	}

}
func process(ch chan []byte) (int64, int64) {
	ok := false
	rb := int64(0)
	tb := int64(0)
	for line := range ch {
		if bytes.Contains(line, []byte(":27018")) {
			ok = true
			continue
		}
		if ok {
			fields := bytes.Split(line, []byte(","))
			for _, field := range fields {
				if bytes.HasPrefix(field, []byte("rb")) {
					size, _ := strconv.ParseInt(field[2:])
					rb += size
					continue
				}
				if bytes.HasPrefix(field, []byte("tb")) {
					size, _ := strconv.ParseInt(field[2:])
					tb += size
					continue
				}
			}
			ok = false
		}
	}
	return rb, tb
}
