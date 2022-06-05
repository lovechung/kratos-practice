package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/lovechung/go-kit/bootstrap"
	"kratos-practice/internal/conf"
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

func newApp(logger log.Logger, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(Service.GetInstanceId()),
		kratos.Name(Service.Name),
		kratos.Version(Service.Version),
		kratos.Metadata(Service.Metadata),
		kratos.Logger(logger),
		kratos.Server(gs),
		kratos.Registrar(rr),
	)
}

// 加载启动配置
func loadConfig() (*conf.Bootstrap, *conf.Registry) {
	c := bootstrap.NewConfigProvider(Flags.Conf,
		Flags.ConfigType,
		Flags.ConfigHost,
		Flags.ConfigToken,
		"conf/"+Service.Name+"/"+Flags.Env+"/data")
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	return &bc, &rc
}

// 加载otel配置
func loadOtel(bc *conf.Bootstrap) {
	shutdownTrace := bootstrap.NewTracerProvider(bc.Otel.Endpoint, Flags.Env, &Service)
	defer shutdownTrace()
	shutdownMetric := bootstrap.NewMetricProvider(bc.Otel.Endpoint, Flags.Env, &Service, false)
	defer shutdownMetric()
}

// 加载日志配置
func loadLogger(bc *conf.Bootstrap) log.Logger {
	return bootstrap.NewLoggerProvider(Flags.Env, bc.Log.File, &Service)
}

func main() {
	flag.Parse()

	bc, rc := loadConfig()
	if bc == nil || rc == nil {
		panic("load config failed")
	}

	loadOtel(bc)

	logger := loadLogger(bc)

	app, cleanup, err := wireApp(bc.Server, bc.Data, rc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
