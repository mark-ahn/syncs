package syncs

import "github.com/mark-ahn/metrics"

type MetricData interface {
	metric_type()
}

type Scope = metrics.Scope[MetricData]

type ThreadCountMetric struct {
	Delta int
}

func (ThreadCountMetric) metric_type() {}
