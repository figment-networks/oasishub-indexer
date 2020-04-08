package report

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type Model struct {
	*shared.Model

	StartHeight  types.Height
	EndHeight    types.Height
	SuccessCount *int64
	ErrorCount   *int64
	ErrorMsg     *string
	Duration     *int64
	Details      types.Jsonb
	CompletedAt  *time.Time
}

// - METHODS
func (Model) TableName() string {
	return "reports"
}

func (r *Model) ValidOwn() bool {
	return r.StartHeight.Valid() &&
		r.EndHeight.Valid()
}

func (r *Model) EqualOwn(m Model) bool {
	return true
}

func (r *Model) Valid() bool {
	return r.Model.Valid() &&
		r.ValidOwn()
}

func (r *Model) Equal(m Model) bool {
	return r.Model != nil &&
		m.Model != nil &&
		r.Model.Equal(*m.Model) &&
		r.EqualOwn(m)
}

func (r *Model) Complete(successCount int64, errorCount int64, err *string, details []byte) {
	completedAt := time.Now()
	duration := completedAt.Sub(r.CreatedAt).Milliseconds()

	r.SuccessCount = &successCount
	r.ErrorCount = &errorCount
	r.ErrorMsg = err
	r.Details = types.Jsonb{RawMessage: details}
	r.Duration = &duration
	r.CompletedAt = &completedAt
}
