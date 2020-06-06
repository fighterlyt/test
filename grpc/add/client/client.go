package main

import (
	"context"
	"fmt"
	"github.com/fighterlyt/test/grpc/add"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
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
					argument := &add.Data{Value: 1}
					if i%2 == 0 {
						argument.Foo = &add.Data_Name{Name: "1"}
					} else {
						argument.Foo = &add.Data_Key{Key: "2"}
					}
					if _, err := client.Add(context.TODO(), argument); err != nil {
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
