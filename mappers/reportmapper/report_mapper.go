package reportmapper

import (
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/domain/reportdomain"
	"github.com/figment-networks/oasishub/utils/errors"
	"github.com/figment-networks/oasishub/utils/log"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func FromPersistence(o orm.ReportModel) (*reportdomain.Report, errors.ApplicationError) {
	details, err := o.Details.MarshalJSON()
	if err != nil {
		log.Error(err)
		return nil, errors.NewError("error when trying to marshal report details", errors.UnmarshalError, err)
	}

	e := &reportdomain.Report{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),

		StartHeight:  o.StartHeight,
		EndHeight:    o.EndHeight,
		SuccessCount: o.SuccessCount,
		ErrorCount:   o.ErrorCount,
		ErrorMsg:     o.ErrorMsg,
		Duration:     o.Duration,
		Details:      details,
		CompletedAt:  o.CompletedAt,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("report not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(r *reportdomain.Report) (*orm.ReportModel, errors.ApplicationError) {
	details := postgres.Jsonb{RawMessage: r.Details}

	if !r.Valid() {
		return nil, errors.NewErrorFromMessage("report not valid", errors.NotValid)
	}

	return &orm.ReportModel{
		EntityModel:  orm.EntityModel{ID: r.ID},
		StartHeight:  r.StartHeight,
		EndHeight:    r.EndHeight,
		SuccessCount: r.SuccessCount,
		ErrorCount:   r.ErrorCount,
		ErrorMsg:     r.ErrorMsg,
		Duration:     r.Duration,
		Details:      details,
		CompletedAt:  r.CompletedAt,
	}, nil
}
