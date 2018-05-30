package main

import (
	"github.com/VividCortex/gohistogram"
	"fmt"
	"github.com/aybabtme/uniplot/histogram"
	"os"
	fighterlyt "github.com/fighterlyt/uniplot/histogram"
)

func main(){
	gram:=gohistogram.NewHistogram(20)
	result:=make([]float64,0,100)
	another:=fighterlyt.BucketsHistro("test",[]float64{float64(10),float64(20),50})
	for i:=0;i<100;i++{
		gram.Add(float64(i))
		another.Add(float64(i))
		result=append(result,float64(i))
	}

	uniGram:=histogram.Hist(20,result)
	histogram.Fprint(os.Stdout,uniGram,histogram.Linear(5))
	fmt.Println(gram.String(),gram.Mean(),gram.Quantile(0.99))
	fighterlyt.Fprint(os.Stdout,*another,fighterlyt.Linear(5))
	another.Generate("test")
}
