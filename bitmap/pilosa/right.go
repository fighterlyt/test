package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pilosa/pilosa/roaring"
)

func main(){
	b:=roaring.NewBitmap()
	count:=10000
	source:=make(map[uint64]struct{},count)
	//r:=rand.New(rand.NewSource(time.Now().UnixNano()))

	for i:=0;i<count;i++{
		temp:=uint64(2*i+1)
		//temp:=uint64(r.Int63n(int64(count*100)))
		//for _,exist:=source[temp];exist;{
		//	println("重复",len(source),i,temp)
		//	panic("1")
		//	//temp=uint64(r.Int63n(int64(count*100)))
		//	temp=r.Uint64()
		//}
		source[temp]=struct{}{}
		b.Add(temp)
	}
	jsonBuffer:=&bytes.Buffer{}

	if err:=json.NewEncoder(jsonBuffer).Encode(source);err!=nil{
		panic(err.Error())
	}
	buffer:=&bytes.Buffer{}

	if _,err:=b.WriteTo(buffer);err!=nil {
		panic(err.Error())
	}
	//}else{
	//	if n!=int64(count){
	//		panic(fmt.Sprintf("数量错误，应该是%d,实际为%d",count,n))
	//	}
	//}

	another:=roaring.NewBitmap()
	if err:=another.UnmarshalBinary(buffer.Bytes());err!=nil{
		panic("反序列化失败"+err.Error())
	}
	//if another.Count()!=uint64(count){
	//	panic(fmt.Sprintf("数量错误，应该是%d,实际为%d",count,another.Count()))
	//}

	fmt.Printf("字节数量%d,序列化字节数量%d,最大值%d,json字节%d\n",b.Size(),buffer.Len(),another.Max(),jsonBuffer.Len())
	println(buffer.Bytes())
	for k:=range source{
		if !another.Contains(uint64(k)){
			panic(fmt.Sprintf("%d不存在",k))
		}
		another.Remove(uint64(k))
	}
	if another.Count()!=0{
		panic(fmt.Sprintf("数量错误，应该是%d,实际为%d",count,another.Count()))
	}
}
