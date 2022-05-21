package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"kratos-practice/internal/conf"
	"kratos-practice/internal/pkg/util/bootstrap"
)

// go build -ldflags "-X main.Service.Version=x.y.z"
var (
	Service = bootstrap.NewServiceInfo(
		"kratos-practice",
		"1.0.0",
		"",
	)

	Flags = bootstrap.NewCommandFlags()
)

func init() {
	Flags.Init()

}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(Service.GetInstanceId()),
		kratos.Name(Service.Name),
		kratos.Version(Service.Version),
		kratos.Metadata(Service.Metadata),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

// 加载配置
func loadConfig() *conf.Bootstrap {
	c := bootstrap.NewConfigProvider(Flags.Conf)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}

func main() {
	flag.Parse()

	bc := loadConfig()
	if bc == nil {
		panic("load config failed")
	}

	shutdownMetric := bootstrap.NewMetricProvider(bc.Otel.Endpoint, bc.Server.Profile, &Service)
	defer shutdownMetric()

	shutdownTrace := bootstrap.NewTracerProvider(bc.Otel.Endpoint, bc.Server.Profile, &Service)
	defer shutdownTrace()

	logger := bootstrap.NewLoggerProvider(bc.Server.Profile, bc.Log, &Service)

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
