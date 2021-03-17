package types

import "time"

const (
	IntervalHourly  SummaryInterval = "hour"
	IntervalDaily   SummaryInterval = "day"
	IntervalMonthly SummaryInterval = "month"
)

// SummaryInterval type represents summary interval
type SummaryInterval string

func (s SummaryInterval) Valid() bool {
	return s == IntervalHourly || s == IntervalDaily || s == IntervalMonthly
}

func (s SummaryInterval) Equal(o SummaryInterval) bool {
	return s == o
}

func (s SummaryInterval) ToDuration() (time.Duration, error) {
	if s == IntervalDaily {
		return time.ParseDuration("24h")
	}
	if s == IntervalMonthly {
		return time.ParseDuration("1m")
	}
	return time.ParseDuration("1h")
}
