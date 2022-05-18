package bootstrap

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"google.golang.org/grpc"
	"time"
)

func NewMetricProvider(endpoint string) func() {
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
	)
	ctx := context.Background()
	exp, err := otlpmetric.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to create the collector exporter: %v", err)
	}

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exp,
		),
		controller.WithExporter(exp),
		controller.WithCollectPeriod(time.Second*2),
	)
	global.SetMeterProvider(pusher)
	err = pusher.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start the collector: %v", err)
	}

	err = host.Start()
	if err != nil {
		log.Fatalf("failed to start the host metric: %v", err)
	}

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
