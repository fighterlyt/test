package main

import (
	"github.com/pltr/onering"
	"sync"
	"fmt"
	"time"
	"log"
)

func main() {
	count := uint32(10000000)


	start:=time.Now()
	ring(count)
	log.Printf("%d 个数字，耗时%s\n",count,time.Since(start).String())

	start=time.Now()
	channel(count)
	log.Printf("%d 个数字，耗时%s\n",count,time.Since(start).String())
}

func ring(count uint32){
	queue := onering.New{Size: count}.SPSC()

	closed := &sync.WaitGroup{}
	closed.Add(2)
	go func() {
		for i := 0; uint32(i) < count; i++ {
			src := int64(i)
			queue.Put(&src)
		}
		closed.Done()
	}()

	go func() {
		var dst *int64
		for i := 0; uint32(i) < count; i++ {
			queue.Get(&dst)
			if *dst != int64(i) {
				panic(fmt.Sprintf("%d/%d", *dst, i))
			}
		}
		closed.Done()
	}()
	closed.Wait()
}
func channel(count uint32){
	closed := &sync.WaitGroup{}
	closed.Add(2)
	ch:=make(chan int64,count)
	go func() {
		for i := 0; uint32(i) < count; i++ {
			src := int64(i)
			ch<-src
		}
		closed.Done()
	}()

	go func() {

		for i := 0; uint32(i) < count; i++ {
			dst:=<-ch
			if dst != int64(i) {
				panic(fmt.Sprintf("%d/%d", dst, i))
			}
		}
		closed.Done()
	}()
	closed.Wait()

}