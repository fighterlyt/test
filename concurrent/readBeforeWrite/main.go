package main

import "sync/atomic"

func main() {
	ch := make(chan int, 1000)
	concurrent := 10
	signal := make(chan int, concurrent)
	wait := new(int64)

	for i := 0; i < concurrent; i++ {
		atomic.AddInt64(wait, 1)
		go read(ch, signal)
	}
	for range signal {
		if atomic.AddInt64(wait, -1) == 0 {
			break
		}
	}
	go write(ch)
}

func read(ch <-chan int, signal chan<- int) {
	select {
	case <-ch:
	case signal <- 1:
	}
}

func write(ch <-chan int) {

}
