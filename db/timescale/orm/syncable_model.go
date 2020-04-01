package orm

import (
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/types"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type SyncableModel struct {
	EntityModel
	SequenceModel

	Type        syncabledomain.Type
	Report      ReportModel `gorm:"foreignkey"`
	ReportID    *types.UUID
	Data        postgres.Jsonb
	ProcessedAt *time.Time
}

func (SyncableModel) TableName() string {
	return "syncables"
}