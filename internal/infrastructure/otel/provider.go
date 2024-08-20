package otel

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Providers contains both the TracerProvider and the MeterProvider
type Providers struct {
	TracerProvider *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
}

// NewProviders initializes both a trace and a metrics provider
func NewProviders(serviceName string) (*Providers, func(), error) {
	// Create a new Prometheus exporter
	promExporter, err := prometheus.New(prometheus.WithNamespace("example"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the Prometheus exporter: %w", err)
	}

	ctx := context.Background()

	otlpExporter, err := otlptracehttp.New(ctx,
		// send this to a local otel collector for now, allow overrides later
		otlptracehttp.WithEndpoint("http://localhost:4318"),
		// we are presuming this happens over a service mesh or locally on the node etc
        otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the otlp exporter: %w", err)
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.DeploymentEnvironmentKey.String("production"),
	)

	// Set up the TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithResource(resources),
		trace.WithBatcher(otlpExporter),
		trace.WithSampler(trace.AlwaysSample()), // Adjust sampler as needed
	)

	// Set up the MeterProvider
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resources),
		metric.WithReader(promExporter),
	)

	// Set the global providers
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(meterProvider)

	// Define the cleanup function
	cleanup := func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down TracerProvider: %v", err)
		}
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down MeterProvider: %v", err)
		}
	}

	return &Providers{
		TracerProvider: tp,
		MeterProvider:  meterProvider,
	}, cleanup, nil
}
