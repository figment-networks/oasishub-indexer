package shared

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type Aggregate struct {
	StartedAtHeight types.Height `json:"started_at_height"`
	StartedAt       types.Time    `json:"started_at"`
}

func (a *Aggregate) Valid() bool {
	return a.StartedAtHeight.Valid() &&
		!a.StartedAt.IsZero()
}

func (a *Aggregate) Equal(m Aggregate) bool {
	return a.StartedAtHeight.Equal(m.StartedAtHeight) &&
		a.StartedAt.Equal(m.StartedAt)
}