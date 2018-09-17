package main

import (
	"log"
	"reflect"
	"sync"
	"time"
)

func main() {
	chanCount := 4
	count := 100000
	chans := make([]chan interface{}, 0, chanCount)
	for i := 0; i < chanCount; i++ {
		chans = append(chans, make(chan interface{}))
	}
	SelectCase(chans, count)
	OriginSelect(chans,count)
}

func SelectCase(chans []chan interface{}, count int) {

	cases := make([]reflect.SelectCase, 0, len(chans))
	for _, ch := range chans {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan interface{}, wg *sync.WaitGroup) {
			for i := 0; i < count; i++ {
				ch <- i
			}
			wg.Done()
		}(ch, wg)
	}
	start := time.Now()

	for i := 0; i < count*len(chans); i++ {
		reflect.Select(cases)
	}

	wg.Wait()
	log.Println(time.Since(start).String())

}

func OriginSelect(chans []chan interface{},count int){


	wg := &sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan interface{}, wg *sync.WaitGroup) {
			for i := 0; i < count; i++ {
				ch <- i
			}
			wg.Done()
		}(ch, wg)
	}
	start := time.Now()

	for i := 0; i < count*len(chans); i++ {
		select{
		case <-chans[0]:
		case <-chans[1]:
		case <-chans[2]:
		case <-chans[3]:
		//case <-chans[4]:
		//case <-chans[5]:
		//case <-chans[6]:
		//case <-chans[7]:
		//case <-chans[8]:
		//case <-chans[9]:

		}
	}

	wg.Wait()
	log.Println(time.Since(start).String())
}