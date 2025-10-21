package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var Log *zap.Logger

func InitTracer(ctx context.Context, serviceName string) func(context.Context) error {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown
}

func InitLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	Log = logger
	return nil
}

func TraceLogger(ctx context.Context) *zap.Logger {
	if Log == nil {
		_ = InitLogger()
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return Log
	}
	return Log.With(
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("span_id", span.SpanContext().SpanID().String()),
	)
}
