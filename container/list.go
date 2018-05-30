package main

import (
	"container/list"
	"fmt"
)

func main(){
	l:=list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.Remove(&list.Element{
		Value:2,
	})
	fmt.Println(l.Len())
}
