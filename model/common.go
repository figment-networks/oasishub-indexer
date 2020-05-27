package model

import "github.com/figment-networks/oasishub-indexer/types"

type Model struct {
	ID        types.ID   `json:"id"`
	CreatedAt types.Time `json:"created_at"`
	UpdatedAt types.Time `json:"updated_at"`
}

func (e *Model) Valid() bool {
	return true
}

func (e *Model) Equal(m Model) bool {
	return e.ID.Equal(m.ID)
}

type Sequence struct {
	Height int64 `json:"height"`
	Time   types.Time   `json:"time"`
}

func (s *Sequence) Valid() bool {
	return s.Height >= 0 &&
		!s.Time.IsZero()
}

func (s *Sequence) Equal(m Sequence) bool {
	return s.Height == m.Height &&
		s.Time.Equal(m.Time)
}

type Aggregate struct {
	StartedAtHeight int64 `json:"started_at_height"`
	StartedAt       types.Time   `json:"started_at"`
	RecentAtHeight  int64 `json:"recent_at_height"`
	RecentAt        types.Time   `json:"recent_at"`
}

func (a *Aggregate) Valid() bool {
	return a.StartedAtHeight >= 0 &&
		!a.StartedAt.IsZero()
}

func (a *Aggregate) Equal(m Aggregate) bool {
	return a.StartedAtHeight == m.StartedAtHeight &&
		a.StartedAt.Equal(m.StartedAt)
}
