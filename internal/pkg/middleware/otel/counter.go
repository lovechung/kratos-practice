package otel

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"kratos-practice/internal/pkg/middleware/i"
)

var _ i.Counter = (*counter)(nil)

type counter struct {
	cnt syncint64.Counter
	lvs []attribute.KeyValue
}

func NewCounter(c syncint64.Counter) i.Counter {
	return &counter{
		cnt: c,
	}
}

func (c *counter) With(lvs []attribute.KeyValue) i.Counter {
	return &counter{
		cnt: c.cnt,
		lvs: lvs,
	}
}

func (c *counter) Add(ctx context.Context, incr int64) {
	c.cnt.Add(ctx, incr, c.lvs...)
}
