package main
//测试过量放入是否会阻塞
import "github.com/pltr/onering"

func main(){
	queue := onering.New{Size: 10}.SPSC()
	for i:=0;i<15;i++{
		src:=int64(i)
		queue.Put(&src)
	}
}
