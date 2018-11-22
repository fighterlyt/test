package main

import (
	"fmt"
	"github.com/nadoo/convtrad"
	"gopkg.in/iconv.v1"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	URL = "http://www.angelfire.com/hi/hayashi/sm%d_%02d.txt?v=%d"
)

var (
	ranges = map[int][]int{
		1: []int{1, 5},
		//2: []int{21, 55},
		//3: []int{56, 90},
	}
)

func main() {


	for max, each := range ranges {
		for start := each[0]; start <= each[1]; start++ {
			fetch(max, start)
			time.Sleep(time.Second*5)
		}
	}

}

func fetch(big, small int) {
	url := fmt.Sprintf(URL, big, small,time.Now().UnixNano())
	println(url)
	cd, err := iconv.Open("utf-8", "big5")
	if err != nil {
		panic(err.Error())
	}

	if resp, err := http.Get(url); err != nil {
		panic(err.Error())
	} else {
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		str := cd.ConvString(string(data))
		println(str)
		if err != nil {
			panic(err.Error())
		} else {
			if file, err := os.Create(fmt.Sprintf("%2d.txt", small)); err != nil {
				panic(err.Error())
			} else {
				defer file.Close()
				fmt.Fprint(file, convtrad.ToSimp(str))
				file.Close()
			}
		}

	}

}
