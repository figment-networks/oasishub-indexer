package types

const (
	IntervalHourly SummaryInterval = "hour"
	IntervalDaily  SummaryInterval = "day"
)

// SummaryInterval type represents summary interval
type SummaryInterval string

func (s SummaryInterval) Valid() bool {
	return s == IntervalHourly || s == IntervalDaily
}

func (s SummaryInterval) Equal(o SummaryInterval) bool {
	return s == o
}