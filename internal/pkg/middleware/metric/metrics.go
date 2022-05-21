package metric

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/attribute"
	metrics "kratos-practice/internal/pkg/middleware/i"
	"time"
)

// Option is metrics option.
type Option func(*options)

// WithRequests with requests counter.
func WithRequests(c metrics.Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
func WithSeconds(c metrics.Histogram) Option {
	return func(o *options) {
		o.seconds = c
	}
}

type options struct {
	// counter: <client/server>_requests_code_total{kind, operation, code, reason}
	requests metrics.Counter
	// histogram: <client/server>_requests_seconds_bucket{kind, operation}
	seconds metrics.Histogram
}

// Server is middleware server-side metrics.
func Server(opts ...Option) middleware.Middleware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			serverAttribute := []attribute.KeyValue{
				attribute.String("kind", kind),
				attribute.String("operation", operation),
				attribute.Int("code", code),
				attribute.String("reason", reason),
			}
			if op.requests != nil {
				op.requests.With(serverAttribute).Add(ctx, 1)
			}
			if op.seconds != nil {
				op.seconds.With(serverAttribute).Record(ctx, time.Since(startTime).Milliseconds())
			}
			return reply, err
		}
	}
}

// Client is middleware client-side metrics.
//func Client(opts ...Option) middleware.Middleware {
//	op := options{}
//	for _, o := range opts {
//		o(&op)
//	}
//	return func(handler middleware.Handler) middleware.Handler {
//		return func(ctx context.Context, req interface{}) (interface{}, error) {
//			var (
//				code      int
//				reason    string
//				kind      string
//				operation string
//			)
//			startTime := time.Now()
//			if info, ok := transport.FromClientContext(ctx); ok {
//				kind = info.Kind().String()
//				operation = info.Operation()
//			}
//			reply, err := handler(ctx, req)
//			if se := errors.FromError(err); se != nil {
//				code = int(se.Code)
//				reason = se.Reason
//			}
//			if op.requests != nil {
//				op.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
//			}
//			if op.seconds != nil {
//				op.seconds.With(kind, operation).Observe(time.Since(startTime).Seconds())
//			}
//			return reply, err
//		}
//	}
//}
