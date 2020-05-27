package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type Report struct {
	*Model

	StartHeight  int64
	EndHeight    int64
	SuccessCount *int64
	ErrorCount   *int64
	ErrorMsg     *string
	Duration     time.Duration
	Details      types.Jsonb
	CompletedAt  *types.Time
}

// - METHODS
func (Report) TableName() string {
	return "reports"
}

func (r *Report) Valid() bool {
	return r.StartHeight >= 0 &&
		r.EndHeight >= 0
}

func (r *Report) Equal(m Report) bool {
	return m.Model != nil &&
		r.Model.ID == m.Model.ID
}

func (r *Report) Complete(successCount int64, errorCount int64, err error) {
	completedAt := types.NewTimeFromTime(time.Now())

	r.SuccessCount = &successCount
	r.ErrorCount = &errorCount
	r.Duration = time.Since(r.CreatedAt.Time)
	r.CompletedAt = completedAt
	//TODO: Implement details
	//r.Details = types.Jsonb{RawMessage: details}

	if err != nil {
		errMsg := err.Error()
		r.ErrorMsg = &errMsg
	}
}
