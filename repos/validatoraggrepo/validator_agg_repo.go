package validatoraggrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/validatoragg"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.PublicKey) bool
	Count() (*int64, errors.ApplicationError)
	GetByEntityUID(types.PublicKey) (*validatoragg.Model, errors.ApplicationError)
	GetAllForHeightGreaterThan(types.Height) ([]validatoragg.Model, errors.ApplicationError)

	// Commands
	Create(*validatoragg.Model) errors.ApplicationError
	Save(*validatoragg.Model) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) *dbRepo {
	return &dbRepo{
		client: c,
	}
}

func (r *dbRepo) Exists(key types.PublicKey) bool {
	q := validatoragg.Model{
		EntityUID: key,
	}
	m := validatoragg.Model{}

	if err := r.client.Where(&q).Take(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(validatoragg.Model{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count for entity aggregate", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of entity aggregate", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByEntityUID(key types.PublicKey) (*validatoragg.Model, errors.ApplicationError) {
	q := validatoragg.Model{
		EntityUID: key,
	}
	var m validatoragg.Model

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find entity aggregate with key %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting validator aggregate by entity UID", errors.QueryError, err)
	}
	return &m, nil
}

func (r *dbRepo) GetAllForHeightGreaterThan(height types.Height) ([]validatoragg.Model, errors.ApplicationError) {
	var ms []validatoragg.Model

	if err := r.client.Where("recent_as_validator_height >= ?", height).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find validator aggregates for height greater than %d", height), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting entity aggregate by entity UID", errors.QueryError, err)
	}
	return ms, nil
}

func (r *dbRepo) Save(m *validatoragg.Model) errors.ApplicationError {
	if err := r.client.Save(m).Error; err != nil {
		msg := "could not save entity aggregate"
		return errors.NewError(msg, errors.CreateError, err)
	}
	return nil
}

func (r *dbRepo) Create(m *validatoragg.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		msg := "could not create entity aggregate"
		return errors.NewError(msg, errors.CreateError, err)
	}
	return nil
}

