package main

import (
	"google.golang.org/grpc"
	"github.com/fighterlyt/test/grpc/add"
	"time"
	"fmt"
	"log"
	"context"
	"sync"
)

func main() {
	if conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure()); err != nil {
		panic(err.Error())
	} else {
		defer conn.Close()
		client := add.NewAddClient(conn)
		start := time.Now()
		startValue := 0
		endValue := 100000
		wg := &sync.WaitGroup{}
		wg.Add(8)
		for i := 0; i < 8; i++ {
			go func() {
				for i := startValue; i < endValue; i++ {
					if _, err := client.Add(context.TODO(), &add.Data{Value: 1}); err != nil {
						panic(fmt.Sprintf("%d-%d,耗时%s,错误:%s", startValue, i, time.Since(start).String(), err.Error()))
					}
				}
				log.Printf("%d-%d,耗时%s", startValue, endValue, time.Since(start).String())
				wg.Done()
			}()
		}

		wg.Wait()
		log.Printf("%d-%d,耗时%s", startValue, endValue*8, time.Since(start).String())

	}
}
