package reportdomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
	"time"
)

type Report struct {
	*commons.DomainEntity

	StartHeight  types.Height
	EndHeight    types.Height
	SuccessCount *int64
	ErrorCount   *int64
	ErrorMsg     *string
	Duration     *int64
	Details      []byte
	CompletedAt  *time.Time
}

// - METHODS
func (r *Report) ValidOwn() bool {
	return r.StartHeight.Valid() &&
		r.EndHeight.Valid()
}

func (r *Report) EqualOwn(m Report) bool {
	return true
}

func (r *Report) Valid() bool {
	return r.DomainEntity.Valid() &&
		r.ValidOwn()
}

func (r *Report) Equal(m Report) bool {
	return r.DomainEntity.Equal(*m.DomainEntity) &&
		r.EqualOwn(m)
}

func (r *Report) Complete(successCount int64, errorCount int64, err *string, details []byte) {
	completedAt := time.Now()
	duration := completedAt.Sub(r.CreatedAt).Milliseconds()

	r.SuccessCount = &successCount
	r.ErrorCount = &errorCount
	r.ErrorMsg = err
	r.Details = details
	r.Duration = &duration
	r.CompletedAt = &completedAt
}
