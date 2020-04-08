package syncable

import (
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type Model struct {
	*shared.Model
	*shared.Sequence

	Type        Type
	Report      report.Model `gorm:"foreignkey"`
	ReportID    *types.ID
	Data        types.Jsonb
	ProcessedAt *time.Time
}

// - Methods
func (Model) TableName() string {
	return "syncables"
}

func (s *Model) ValidOwn() bool {
	return s.Type.Valid()
}

func (s *Model) EqualOwn(m Model) bool {
	return true
}

func (s *Model) Valid() bool {
	return s.Model.Valid() &&
		s.Sequence.Valid() &&
		s.ValidOwn()
}

func (s *Model) Equal(m Model) bool {
	return s.Model != nil &&
		m.Model != nil &&
		s.Model.Equal(*m.Model) &&
		s.Sequence.Equal(*m.Sequence) &&
		s.EqualOwn(m)
}

func (s *Model) MarkProcessed(reportID types.ID) {
	t := time.Now()
	rid := reportID

	s.ProcessedAt = &t
	s.ReportID = &rid
}
