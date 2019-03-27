package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	roaring2 "github.com/RoaringBitmap/roaring"
	"github.com/pilosa/pilosa/roaring"
	"math/rand"
	"time"
)

func main() {
	b := roaring.NewBitmap()
	r := roaring2.NewBitmap()
	count := 10000
	max:=int64(123456789)
	source := make(map[uint64]struct{}, count)
	ran:=rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		temp := uint64(ran.Int63n(max)+1)
		for _,exist:=source[temp];exist;{
			temp = uint64(ran.Int63n(max)+1)
		}
		source[temp] = struct{}{}
	}

	start:=time.Now()
	for value:=range source{
		b.Add(value)
	}
	fmt.Printf("pilosa耗时%s\n",time.Since(start).String())

	start=time.Now()
	for value:=range source{
		r.Add(uint32(value))
	}
	fmt.Printf("std耗时%s\n",time.Since(start).String())

	jsonBuffer := &bytes.Buffer{}
	if err := json.NewEncoder(jsonBuffer).Encode(source); err != nil {
		panic(err.Error())
	}
	buffer := &bytes.Buffer{}

	start=time.Now()
	if _, err := b.WriteTo(buffer); err != nil {
		panic(err.Error())
	}
	fmt.Printf("pilosa序列化耗时%s\n",time.Since(start).String())

	start=time.Now()

	binary,err:=r.MarshalBinary()
	if err!=nil{
		panic(err.Error())
	}
	fmt.Printf("std序列化耗时%s\n",time.Since(start).String())

	fmt.Printf("字节数量%d,序列化字节数量%d,%f%%,最大值%d,json字节%d,std 大小%d %f%%\n", b.Size(), buffer.Len(),div(buffer.Len()*100,jsonBuffer.Len()), b.Max(), jsonBuffer.Len(),len(binary),div(len(binary)*100,jsonBuffer.Len()))

}

func div(a,b int) float64{
	return float64(a)/float64(b)
}