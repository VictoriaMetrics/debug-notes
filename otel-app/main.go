package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var meter = otel.Meter("go.opentelemetry.io/contrib/examples/otel-collector")

type metrics struct {
	requests    metric.Int64Counter
	duration    metric.Float64Histogram
	durationExp metric.Float64Histogram // enhanced exponential
}

var m metrics

func main() {

	ctx := context.Background()
	shutdownMeterProvider := initMetricProvider(ctx)
	defer func() {
		if err := shutdownMeterProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown MeterProvider: %s", err)
		}
	}()

	var err error
	m.requests, err = meter.Int64Counter("service.requests.total", metric.WithDescription("The number of served requests"))
	if err != nil {
		log.Fatal(err)
	}

	m.duration, err = meter.Float64Histogram(
		"request.duration",
		metric.WithDescription("a histogram with custom buckets and rename"),
	)
	if err != nil {
		log.Fatal(err)
	}

	m.durationExp, err = meter.Float64Histogram(
		"request.duration2",
		metric.WithDescription("a histogram with custom buckets and rename"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", handle(home))
	http.HandleFunc("/articles", handle(articles))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		attrs := metric.WithAttributes(attribute.String("path", r.URL.Path))
		m.requests.Add(context.Background(), 1, attrs)
		h.ServeHTTP(w, r)
		m.duration.Record(context.Background(), time.Since(start).Seconds(), attrs)
		m.durationExp.Record(context.Background(), time.Since(start).Seconds(), attrs)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!")
}

func articles(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultClient.Get("https://duckduckgo.com/?t=h_&q=observability")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(data))
}

func initMetricProvider(ctx context.Context) func(context.Context) error {
	conn, err := initConn()
	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String("test-service"),
			semconv.ServiceInstanceID("test-service"),
			semconv.ServerAddressKey.String("test-service"),
			semconv.ServiceNamespace("test-service"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create metrics exporter: %s", err)
	}

	stdoutMetricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		log.Fatalf("Failed to create stdout metric exporter: %v", err)
	}

	exponentialHistogramView := sdkmetric.NewView(
		// The instrument to which this view is attached.
		sdkmetric.Instrument{
			Name: "request.duration2",
			Kind: sdkmetric.InstrumentKindHistogram,
		},
		// The aggregation to use for this view.
		sdkmetric.Stream{
			Name: "request.exponential.duration",
			Aggregation: sdkmetric.AggregationBase2ExponentialHistogram{
				MaxSize:  160, // Maximum number of buckets.
				MaxScale: 20,  // Maximum resolution scale.
			},
		},
	)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(1*time.Second))),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(stdoutMetricExporter, sdkmetric.WithInterval(1*time.Second))),
		sdkmetric.WithResource(res),
		sdkmetric.WithView(exponentialHistogramView),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown
}

// Initialize a gRPC connection to be used by both the tracer and meter
// providers.
func initConn() (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	// You may replace `localhost:4317` with your endpoint.
	conn, err := grpc.NewClient("localhost:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}
