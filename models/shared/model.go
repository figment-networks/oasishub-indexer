package shared

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type Model struct {
	ID        types.ID  `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//- Methods
func (e *Model) Valid() bool {
	return e == nil
}

func (e *Model) Equal(m Model) bool {
	return e.ID.Equal(m.ID)
}
