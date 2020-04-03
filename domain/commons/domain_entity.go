package commons

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type DomainEntity struct {
	ID        types.UUID `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	IsNew bool `json:"-"`
}

type EntityProps struct {
	ID types.UUID
}

func NewDomainEntity(props EntityProps) *DomainEntity {
	id := props.ID
	isNew := false
	if !props.ID.Valid() {
		id = types.NewUUID()
		isNew = true
	}

	return &DomainEntity{
		ID:        id,
		IsNew:     isNew,
		CreatedAt: time.Now(),
	}
}

//- Methods
func (e *DomainEntity) Valid() bool {
	return e.ID.Valid()
}

func (e *DomainEntity) Equal(m DomainEntity) bool {
	return e.ID.Equal(m.ID)
}
