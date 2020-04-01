package accountaggrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/accountdomain"
	"github.com/figment-networks/oasishub-indexer/mappers/accountaggmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.PublicKey) bool
	Count() (*int64, errors.ApplicationError)
	GetByPublicKey(types.PublicKey) (*accountdomain.AccountAgg, errors.ApplicationError)

	// Commands
	Create(*accountdomain.AccountAgg) errors.ApplicationError
	Save(*accountdomain.AccountAgg) errors.ApplicationError
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
	query := orm.AccountAggModel{
		PublicKey: key,
	}
	foundSyncableValidator := orm.AccountAggModel{}

	if err := r.client.Where(&query).Take(&foundSyncableValidator).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(orm.AccountAggModel{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("account aggregate not found", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of account aggregate", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByPublicKey(key types.PublicKey) (*accountdomain.AccountAgg, errors.ApplicationError) {
	query := orm.AccountAggModel{
		PublicKey: key,
	}
	var seq orm.AccountAggModel

	if err := r.client.Where(&query).Take(&seq).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find account aggregate with key %s", key), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting account aggregate by public key", errors.QueryError, err)
	}
	e, err := accountaggmapper.FromPersistence(seq)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *dbRepo) Save(sv *accountdomain.AccountAgg) errors.ApplicationError {
	pr, err := accountaggmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not save account aggregate", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(sv *accountdomain.AccountAgg) errors.ApplicationError {
	b, err := accountaggmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create account aggregate", errors.CreateError, err)
	}
	return nil
}
