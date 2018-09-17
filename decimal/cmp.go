package main

import (
	"math/rand"
	"github.com/fighterlyt/decimal"
	"time"
	"log"
)

func main(){
	count:=10000
	start:=time.Now()
	cmp(count)
	log.Printf("比较%d对，耗时%s\n",count,time.Since(start).String())

}

func cmp(count int){
	r:=rand.New(rand.NewSource(1))
	for i:=0;i<count;i++{
		one:=decimal.NewFromFloat(r.Float64())
		another:=decimal.NewFromFloat(r.Float64())
		one.Cmp(another)
	}
}