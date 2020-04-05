package entityaggrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/entityagg"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.PublicKey) bool
	Count() (*int64, errors.ApplicationError)
	GetByEntityUID(types.PublicKey) (*entityagg.Model, errors.ApplicationError)

	// Commands
	Create(*entityagg.Model) errors.ApplicationError
	Save(*entityagg.Model) errors.ApplicationError
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
	q := entityagg.Model{
		EntityUID: key,
	}
	m := entityagg.Model{}

	if err := r.client.Where(&q).Take(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(entityagg.Model{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count for entity aggregate", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of entity aggregate", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByEntityUID(key types.PublicKey) (*entityagg.Model, errors.ApplicationError) {
	q := entityagg.Model{
		EntityUID: key,
	}
	var m entityagg.Model

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find entity aggregate with key %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting entity aggregate by entity UID", errors.QueryError, err)
	}
	return &m, nil
}

func (r *dbRepo) Save(m *entityagg.Model) errors.ApplicationError {
	if err := r.client.Save(m).Error; err != nil {
		msg := "could not save entity aggregate"
		return errors.NewError(msg, errors.CreateError, err)
	}
	return nil
}

func (r *dbRepo) Create(m *entityagg.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		msg := "could not create entity aggregate"
		return errors.NewError(msg, errors.CreateError, err)
	}
	return nil
}

