package otelutil

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
)

// SetupOTelSDK bootstraps the OpenTelemetry pipeline.
// It takes the OTLP endpoint and tags (comma-separated key:value pairs) as input.
// It returns shutdown function that should be called for proper cleanup.
func SetupOTelSDK(ctx context.Context, epurl, tags string) (func(context.Context) error, error) {
	var (
		shutdownFuncs []func(ctx context.Context) error
		shutdown      = func(ctx context.Context) error {
			var err error
			for _, fn := range shutdownFuncs {
				err = errors.Join(err, fn(ctx))
			}
			shutdownFuncs = nil
			return err
		}
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	var args []attribute.KeyValue
	for i, p := range strings.Split(tags, ",") {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag %q at index %d", kv, i)
		}

		args = append(
			args,
			attribute.String(
				strings.ToValidUTF8(kv[0], "�"),
				strings.ToValidUTF8(kv[1], "�"),
			),
		)
	}

	val, err := url.Parse(epurl)
	if err != nil {
		return nil, err
	}

	var client otlptrace.Client
	switch val.Scheme {
	case "https":
		client = otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(val.Host),
		)
	case "http":
		client = otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(val.Host),
		)
	case "grpc":
		// nolint:staticcheck
		client = otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(val.Host),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)
	default:
		return nil, fmt.Errorf("unsupported scheme %q in: %s", val.Scheme, val)
	}

	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, exp.Shutdown)

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			args...,
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	otel.SetTracerProvider(tracerProvider)

	return shutdown, nil
}
