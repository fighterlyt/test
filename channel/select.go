package main

import "fmt"

func main() {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	go func(){

		close(ch2)
	}()

	for {
		select {
		case value, ok := <-ch1:
			if !ok{
				fmt.Println("ch1")
			}else{
				fmt.Println(value)
			}
		case value, ok := <-ch2:
			if !ok {
				fmt.Println("ch2")
			}else{
				fmt.Println(value)
			}
		}
	}
}
