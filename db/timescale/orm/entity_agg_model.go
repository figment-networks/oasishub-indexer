package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type EntityAggModel struct {
	EntityModel
	AggregateModel

	EntityUID types.PublicKey
}

func (EntityAggModel) TableName() string {
	return "entity_aggregates"
}
