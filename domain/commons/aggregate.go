package commons

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type Aggregate struct {
	StartedAtHeight types.Height `json:"started_at_height"`
	StartedAt       time.Time    `json:"started_at"`
}

type AggregateProps struct {
	StartedAtHeight types.Height
	StartedAt       time.Time
}

func NewAggregate(props AggregateProps) *Aggregate {
	return &Aggregate{
		StartedAtHeight: props.StartedAtHeight,
		StartedAt:       props.StartedAt,
	}
}

func (a *Aggregate) GetStartedAtHeight() types.Height { return a.StartedAtHeight }
func (a *Aggregate) GetStartedAt() time.Time          { return a.StartedAt }

func (a *Aggregate) Valid() bool {
	return a.StartedAtHeight.Valid() &&
		!a.StartedAt.IsZero()
}

func (a *Aggregate) Equal(m Aggregate) bool {
	return a.StartedAtHeight.Equal(m.StartedAtHeight) &&
		a.StartedAt.Equal(m.StartedAt)
}
