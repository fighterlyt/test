package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"log"
	"net/http"
	"time"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func config() (io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, err
	}
	return closer, nil
}
func init() {
	var err error
	if closer, err = config(); err != nil {
		log.Fatal(err.Error())
	} else {
		tracer = opentracing.GlobalTracer()
	}
}

var tracer = opentracing.GlobalTracer()
var closer io.Closer

func main() {

	defer closer.Close()

	http.HandleFunc("/1", first)
	http.HandleFunc("/2", second)
	http.ListenAndServe(":1234", nil)

}

func first(resp http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("first")

	second, _ := http.NewRequest(http.MethodGet, "http://localhost:1234/2", nil)
	carrier := opentracing.HTTPHeadersCarrier(second.Header)

	tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	http.DefaultClient.Do(second)
	defer span.Finish()
	println("first")
}

func second(resp http.ResponseWriter, req *http.Request) {
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	sc, _ := tracer.Extract(opentracing.HTTPHeaders, carrier)
	span := tracer.StartSpan("second", opentracing.ChildOf(sc))
	time.Sleep(time.Millisecond * 10)
	span.Finish()
	println("second")
}
