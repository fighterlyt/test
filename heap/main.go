package main

import (
	"fmt"
	"runtime"
)

func main() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	fmt.Println(m.Alloc)
}
