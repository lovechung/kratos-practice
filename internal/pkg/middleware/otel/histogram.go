package otel

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"kratos-practice/internal/pkg/middleware/i"
)

var _ i.Histogram = (*histogram)(nil)

type histogram struct {
	his syncint64.Histogram
	lvs []attribute.KeyValue
}

// NewHistogram new a prometheus histogram and returns Histogram.
func NewHistogram(h syncint64.Histogram) i.Histogram {
	return &histogram{
		his: h,
	}
}

func (h histogram) With(lvs []attribute.KeyValue) i.Histogram {
	return &histogram{
		his: h.his,
		lvs: lvs,
	}
}

func (h histogram) Record(ctx context.Context, incr int64) {
	h.his.Record(ctx, incr, h.lvs...)
}
