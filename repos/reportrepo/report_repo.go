package reportrepo

import (
	"github.com/figment-networks/oasishub/domain/reportdomain"
	"github.com/figment-networks/oasishub/mappers/reportmapper"
	"github.com/figment-networks/oasishub/utils/errors"
	"github.com/figment-networks/oasishub/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Commands
	Create(*reportdomain.Report) errors.ApplicationError
	Save(*reportdomain.Report) errors.ApplicationError
}

type dbRepo struct{
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) DbRepo {
	return &dbRepo{
		client: c,
	}
}

func (r *dbRepo) Create(report *reportdomain.Report) errors.ApplicationError {
	b, err := reportmapper.ToPersistence(report)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create block", errors.CreateError, err)
	}
	return nil
}

func (r *dbRepo) Save(report *reportdomain.Report) errors.ApplicationError {
	pr, err := reportmapper.ToPersistence(report)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create report", errors.CreateError, err)
	}
	return nil
}

