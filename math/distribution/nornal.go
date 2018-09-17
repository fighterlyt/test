package main

import (
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"math"
	"fmt"
)

func main() {
	plt.Reset(false, nil)

	normal(0.0,1.0,1,"0","1.0")
	normal(0.0,0.2,2,"0.0","0.2")

	normal(0.0,5.0,3,"0.0","5.0")

	normal(-2.0,0.5,4,"-2.0","0.5")

	// save figure
	plt.Save("/tmp/gosl", "normal")

}

func normal(μ,σ float64,index int,xName,yName string){
	// initialise generator
	rnd.Init(1234)

	σ=math.Sqrt(σ)

	// generate samples
	nsamples := 100
	tick:=5*2/100
	xmin, xmax := -5.0,5.0// limits for histogram
	X := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		X[i] = xmin+float64(tick*i)
	}

	// constants
	nstations := 41        // number of bins + 1


	// build histogram: count number of samples within each bin
	var hist rnd.Histogram
	hist.Stations = utl.LinSpace(xmin, xmax, nstations)
	hist.Count(X, true)

	// compute area of density diagram
	area := hist.DensityArea(nsamples)
	io.Pf("area = %v\n", area)

	// plot lognormal distribution
	var dist rnd.DistNormal
	dist.Init(&rnd.Variable{M: μ, S: σ})

	// compute lognormal points for plot
	x := utl.LinSpace(xmin, xmax, nsamples)
	y := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		y[i] = dist.Pdf(x[i])

	}

	// plot density
	//plt.Subplot(2, 1, index)
	plt.Plot(x, y, &plt.A{
		L:fmt.Sprintf("%f,%f",μ,σ),
	})
	plt.AxisYmax(1.0)
	plt.Gll("$"+xName+"$", "$"+yName+"$", nil)

}