package main

import "fmt"

func main(){
	count:=10
	index:=-1
	smallest:=float64(count*2)
	for i:=1;i<count;i++{
		amount:=float64(i)+float64(count/i)
		if amount<smallest{
			smallest=amount
			index=i
		}
	}
	fmt.Println("smallest",index,smallest)
}
