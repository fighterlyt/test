package main

import (
	"github.com/cpmech/gosl/plt"
	"fmt"
	"math"
	"github.com/cpmech/gosl/rnd"
)

func main() {
	count := 60
	result := make([]float64, 0, count)
	rate := 0.99
	currentRate:=0.99
	x:=make([]float64,0,count)
	amounts:=make([]float64,0,count)
	for i := 0; i < count; i++ {
		value := 0.01 * currentRate
		currentRate = currentRate * rate
		result=append(result,value)
		x=append(x,float64(i))
		amounts=append(amounts,amount(int(value*1000),value))
	}
	plt.Reset(false, nil)
	plt.Subplot(2, 1, 1)
	plt.Plot(x, result, nil)
	plt.Gll("$x$", "$f(x)$", nil)

	plt.Subplot(2,1,2)
	plt.Plot(x,amounts,nil)
	plt.Gll("$x$", "$amount$", nil)

	// save figure
	plt.Save("/tmp/gosl", "manual")
	for _,value:=range result{
		fmt.Println(value,int(value*1000))
	}
}


func amount(k int,price float64) float64{
	m:=1.0

	v:=1.0

	d:=1.0/1.0001
	ss:=1.0001

	one:=math.Pow(d,math.Abs(float64(k-500)))*m
	another:=math.Pow(ss,math.Abs(float64(k-500)))*v

	dist:=rnd.DistNormal{}
	dist.Init(&rnd.Variable{M:one,S:another})
	return dist.Pdf(price)


}