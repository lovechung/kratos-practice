package bootstrap

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"time"
)

func NewTracerProvider(endpoint, env string, serviceInfo *ServiceInfo) func() {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	ctx := context.Background()
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to create the trace: %v", err)
	}

	bsp := traceSdk.NewBatchSpanProcessor(exp)
	tracerProvider := traceSdk.NewTracerProvider(
		traceSdk.WithSampler(traceSdk.AlwaysSample()),
		traceSdk.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(serviceInfo.Name),
			semConv.ServiceVersionKey.String(serviceInfo.Version),
			semConv.ServiceInstanceIDKey.String(serviceInfo.Id),
			attribute.String("env", env),
		)),
		traceSdk.WithSpanProcessor(bsp),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
