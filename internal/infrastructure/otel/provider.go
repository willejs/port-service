package otel

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
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

	// todo: inject this from config
	jaegerEndpoint := "http://localhost:14268/api/traces"
	jaeggerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the Jaeger exporter: %w", err)
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
		trace.WithBatcher(jaeggerExporter),
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
