package main

import "testing"

func BenchmakrSpeed(b *testing.B){

	count:=uint32(100000)
	ring(count)
}
