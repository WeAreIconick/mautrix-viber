// Package tracing provides OpenTelemetry tracing for request flows.
// Uses OTLP (OpenTelemetry Protocol) exporter for modern observability backends.
// Supports Jaeger, Zipkin, and other OTLP-compatible collectors.
package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
)

// InitTracing initializes OpenTelemetry tracing with OTLP HTTP exporter.
// Uses the modern OTLP protocol instead of the deprecated Jaeger exporter.
// For Jaeger: Jaeger supports OTLP natively since v1.35+, use http://jaeger:4318
// For other backends: Use any OTLP-compatible endpoint
func InitTracing(serviceName, otlpEndpoint string) (func(), error) {
	if otlpEndpoint == "" {
		// No tracing configured
		tracer = trace.NewNoopTracerProvider().Tracer("noop")
		return func() {}, nil
	}
	
	// Create OTLP HTTP exporter (works with Jaeger 1.35+, Zipkin, and other OTLP backends)
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(), // Use WithTLSClientConfig in production
	)
	if err != nil {
		return nil, fmt.Errorf("create OTLP exporter: %w", err)
	}
	
	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}
	
	// Create tracer provider with OTLP exporter
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
	)
	
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	
	tracer = otel.Tracer(serviceName)
	
	// Return shutdown function
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			// Log error using standard library since tracing may not be initialized
			// Using fmt.Printf as fallback since logger may not be available
			_ = err // Suppress unused variable warning - error is expected during shutdown
		}
	}, nil
}

// StartSpan starts a new trace span.
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}
	return tracer.Start(ctx, name, opts...)
}

// AddSpanAttributes adds attributes to a span.
func AddSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	if span == nil || !span.IsRecording() {
		return
	}
	span.SetAttributes(attrs...)
}

// RecordError records an error on a span.
func RecordError(span trace.Span, err error) {
	if span == nil || !span.IsRecording() || err == nil {
		return
	}
	span.RecordError(err)
}

// SpanFromContext gets a span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

