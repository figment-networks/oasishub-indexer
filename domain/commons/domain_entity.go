package commons

import (
	"github.com/figment-networks/oasishub/types"
	"time"
)

type DomainEntity struct {
	ID        types.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	IsNew bool
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
		ID:    id,
		IsNew: isNew,
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
