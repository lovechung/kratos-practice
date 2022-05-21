package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"kratos-practice/api/v1"
	"kratos-practice/internal/conf"
	"kratos-practice/internal/pkg/middleware/otel"
	"kratos-practice/internal/service"

	metrics "kratos-practice/internal/pkg/middleware/metric"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, us *service.UserService, cs *service.CarService, logger log.Logger) *http.Server {
	meter := global.Meter("kratos-practice")
	requestHistogram, _ := meter.SyncInt64().Histogram("request_seconds", instrument.WithUnit(unit.Milliseconds))

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			metrics.Server(
				metrics.WithSeconds(otel.NewHistogram(requestHistogram)),
			),
			logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)

	v1.RegisterUserHTTPServer(srv, us)
	v1.RegisterCarHTTPServer(srv, cs)
	return srv
}
