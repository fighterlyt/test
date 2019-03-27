package main

import (
	"math/big"
	"time"
)
var equal bool
func main(){
	a:=big.NewInt(0)
	b:=big.NewInt(0)
	count:=1000000
	start:=time.Now()
	for i:=0;i<count;i++{
		equal=a.Cmp(b)<=0
	}
	 println(time.Since(start).String())
}
