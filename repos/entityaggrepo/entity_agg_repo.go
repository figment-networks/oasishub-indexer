package entityaggrepo


import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/entitydomain"
	"github.com/figment-networks/oasishub-indexer/mappers/entityaggmapper"
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
	GetByEntityUID(types.PublicKey) (*entitydomain.EntityAgg, errors.ApplicationError)

	// Commands
	Create(*entitydomain.EntityAgg) errors.ApplicationError
	Save(*entitydomain.EntityAgg) errors.ApplicationError
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
	query := orm.EntityAggModel{
		EntityUID: key,
	}
	foundSyncableValidator := orm.EntityAggModel{}

	if err := r.client.Where(&query).Take(&foundSyncableValidator).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(orm.EntityAggModel{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count for entity aggregate", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of entity aggregate", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByEntityUID(key types.PublicKey) (*entitydomain.EntityAgg, errors.ApplicationError) {
	query := orm.EntityAggModel{
		EntityUID: key,
	}
	var seq orm.EntityAggModel

	if err := r.client.Where(&query).Take(&seq).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find entity aggregate with key %s", key), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting entity aggregate by entity UID", errors.QueryError, err)
	}
	e, err := entityaggmapper.FromPersistence(seq)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *dbRepo) Save(sv *entitydomain.EntityAgg) errors.ApplicationError {
	pr, err := entityaggmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		msg := "could not save entity aggregate"
		log.Error(err)
		return errors.NewError(msg, errors.CreateError, err)
	}
	return nil
}

func (r *dbRepo) Create(sv *entitydomain.EntityAgg) errors.ApplicationError {
	b, err := entityaggmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		msg := "could not create entity aggregate"
		log.Error(err)
		return errors.NewError(msg, errors.CreateError, err)
	}
	return nil
}

