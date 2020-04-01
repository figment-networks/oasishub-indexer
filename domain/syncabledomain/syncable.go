package syncabledomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
	"time"
)

type Syncable struct {
	*commons.DomainEntity
	*commons.Sequence

	Type 		 Type
	ReportID     *types.UUID
	Data         []byte
	ProcessedAt  *time.Time
}

func (s *Syncable) ValidOwn() bool {
	return s.Type.Valid()
}

func (s *Syncable) EqualOwn(m Syncable) bool {
	return true
}

func (s *Syncable) Valid() bool {
	return s.DomainEntity.Valid() &&
		s.Sequence.Valid() &&
		s.ValidOwn()
}

func (s *Syncable) Equal(m Syncable) bool {
	return s.DomainEntity.Equal(*m.DomainEntity) &&
		s.Sequence.Equal(*m.Sequence) &&
		s.EqualOwn(m)
}

func (s *Syncable) MarkProcessed(reportID types.UUID) {
	t := time.Now()
	rid := reportID

	s.ProcessedAt = &t
	s.ReportID = &rid
}
