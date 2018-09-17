package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
	"fmt"
)

func main() {

	// initialise generator
	rnd.Init(1234)

	// parameters
	μ := 1.0
	σ := 0.25

	// generate samples
	nsamples := 1000
	X := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		X[i] = rnd.Lognormal(μ, σ)
	}

	// constants
	nstations := 41        // number of bins + 1
	xmin, xmax := 0.0, 3.0 // limits for histogram

	// build histogram: count number of samples within each bin
	var hist rnd.Histogram
	hist.Stations = utl.LinSpace(xmin, xmax, nstations)
	hist.Count(X, true)

	// compute area of density diagram
	area := hist.DensityArea(nsamples)
	io.Pf("area = %v\n", area)

	// plot lognormal distribution
	var dist rnd.DistLogNormal
	dist.Init(&rnd.Variable{M: μ, S: σ})

	// compute lognormal points for plot
	n := 101
	x := utl.LinSpace(0, 3, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = 1-dist.Cdf(x[i])
	}

	// plot density
	plt.Reset(false, nil)
	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, nil)
	hist.PlotDensity(nil)
	plt.AxisYmax(1.8)
	plt.Gll("$x$", "$f(x)$", nil)

	// plot cumulative function
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, nil)
	plt.Gll("$x$", "$F(x)$", nil)

	// save figure
	plt.Save("/tmp/gosl", "rnd_lognormalDistribution")

	for _,value:=range  Y{
		fmt.Println(value)
	}
}