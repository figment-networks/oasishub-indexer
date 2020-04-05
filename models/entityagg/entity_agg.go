package entityagg

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Aggregate

	EntityUID types.PublicKey `json:"entity_uid"`
}

// - METHODS
func (Model) TableName() string {
	return "entity_aggregates"
}

func (aa *Model) ValidOwn() bool {
	return aa.EntityUID.Valid()
}

func (aa *Model) EqualOwn(m Model) bool {
	return aa.EntityUID.Equal(m.EntityUID)
}

func (aa *Model) Valid() bool {
	return aa.Model.Valid() &&
		aa.Aggregate.Valid() &&
		aa.ValidOwn()
}

func (aa *Model) Equal(m Model) bool {
	return aa.Model.Equal(*m.Model) &&
		aa.Aggregate.Equal(*m.Aggregate) &&
		aa.EqualOwn(m)
}

func (aa Model) UpdateAggAttrs(entity Model) {
	// Nothing to update yet
}

