package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type ReportModel struct {
	EntityModel

	StartHeight  types.Height
	EndHeight    types.Height
	SuccessCount *int64
	ErrorCount   *int64
	ErrorMsg     *string
	Duration     *int64
	Details      postgres.Jsonb
	CompletedAt  *time.Time
}

func (ReportModel) TableName() string {
	return "reports"
}
