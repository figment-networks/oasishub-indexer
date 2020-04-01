package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type EntityModel struct {
	ID        types.UUID `gorm:"type:uuid; primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SequenceModel struct {
	ChainId string
	Height  types.Height
	Time    time.Time
}

type AggregateModel struct {
	StartedAtHeight types.Height
	StartedAt       time.Time
}
