package metrics

import (
	"sort"
	"time"
)

type Collector struct {
	metrics Metrics
}

type Metrics struct {
	TotalRequests int
	Errors        int
	Durations     []time.Duration
	StartTime     time.Time
	EndTime       time.Time
}

func NewCollector() *Collector {
	return &Collector{
		metrics: Metrics{
			Durations: make([]time.Duration, 0),
		},
	}
}

func (c *Collector) RecordResult(duration time.Duration, err error) {
	c.metrics.TotalRequests++

	if err != nil {
		c.metrics.Errors++
	} else {
		c.metrics.Durations = append(c.metrics.Durations, duration)
	}
}

func (c *Collector) GetMetrics() Metrics {
	return c.metrics
}

func (c *Collector) CalculatePercentile(p float64) time.Duration {
	sort.Slice(c.metrics.Durations, func(i, j int) bool {
		return c.metrics.Durations[i] < c.metrics.Durations[j]
	})

	index := int(float64(len(c.metrics.Durations)) * p)
	if index >= len(c.metrics.Durations) {
		index = len(c.metrics.Durations) - 1
	}
	return c.metrics.Durations[index]
}

func (c *Collector) AverageDuration() time.Duration {
	if len(c.metrics.Durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range c.metrics.Durations {
		total += d
	}
	return total / time.Duration(len(c.metrics.Durations))
}
