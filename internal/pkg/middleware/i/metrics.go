package i

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
)

// Counter is metrics counter.
type Counter interface {
	With(lvs []attribute.KeyValue) Counter
	Add(ctx context.Context, incr int64)
}

type Histogram interface {
	With(lvs []attribute.KeyValue) Histogram
	Record(ctx context.Context, incr int64)
}
