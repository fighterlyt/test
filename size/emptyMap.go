package main

import (
	"unsafe"
)

func main() {
	m := make(map[string]interface{})
	m["a"]="b"
	println(unsafe.Sizeof(m),unsafe.Sizeof(unsafe.Pointer(&m)))
}
