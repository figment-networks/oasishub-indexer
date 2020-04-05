package accountaggrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/accountagg"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.PublicKey) bool
	Count() (*int64, errors.ApplicationError)
	GetByPublicKey(types.PublicKey) (*accountagg.Model, errors.ApplicationError)

	// Commands
	Create(*accountagg.Model) errors.ApplicationError
	Save(*accountagg.Model) errors.ApplicationError
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
	q := accountagg.Model{
		PublicKey: key,
	}
	m := accountagg.Model{}

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(accountagg.Model{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("account aggregate not found", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of account aggregate", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByPublicKey(key types.PublicKey) (*accountagg.Model, errors.ApplicationError) {
	q := accountagg.Model{
		PublicKey: key,
	}
	var m accountagg.Model

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find account aggregate with key %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting account aggregate by public key", errors.QueryError, err)
	}
	return &m, nil
}

func (r *dbRepo) Save(m *accountagg.Model) errors.ApplicationError {
	if err := r.client.Save(m).Error; err != nil {
		return errors.NewError("could not save account aggregate", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(m *accountagg.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		return errors.NewError("could not create account aggregate", errors.CreateError, err)
	}
	return nil
}
