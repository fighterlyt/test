package main

import (
	"github.com/globalsign/mgo"
	"time"
	"sync"
	"fmt"
	"orderbook/orderbooks/orders"
)

/*	总数量	并发量	耗时
	100000	8		5.6S
	100000	16		5.6S
	100000	30		5.1S
	100000	100		5.1S

 */
func main() {

	count := 100000
	concurrent := 100
	session, _ := mgo.Dial(":27017")

	type A struct {
		B int
	}
	data :=orders.NewOrder()
	start := time.Now()

	step := count / concurrent
	index := 0
	wg := &sync.WaitGroup{}
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {

		index += step
		if index > count {
			step -= (index - count)
		}
		fmt.Println("step", step)
		go func(count int) {
			fmt.Println("count", count)
			for i := 0; i < count; i++ {
				session.DB("test").C("test").Insert(data)

			}
			wg.Done()
		}(step)

	}
	wg.Wait()
	fmt.Println(time.Since(start).String())
}

//func main() {
//
//	count := 100000
//	concurrent := 100
//	server, _ := mmongodb.NewServer("localhost:27017", "", "", "", "")
//	dao := server.NewDAO("test", "test")
//	type A struct {
//		B int
//	}
//	data := orders.NewOrder()
//	start := time.Now()
//
//	step := count / concurrent
//	index := 0
//	wg := &sync.WaitGroup{}
//	wg.Add(concurrent)
//	for i := 0; i < concurrent; i++ {
//
//		index += step
//		if index > count {
//			step -= (index - count)
//		}
//		fmt.Println("step", step)
//		go func(count int) {
//			fmt.Println("count", count)
//			for i := 0; i < count; i++ {
//				dao.Insert(data)
//
//			}
//			wg.Done()
//		}(step)
//
//	}
//	wg.Wait()
//	fmt.Println(time.Since(start).String())
//}
