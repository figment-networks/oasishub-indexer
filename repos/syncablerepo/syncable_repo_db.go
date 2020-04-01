package syncablerepo

import (
	"fmt"
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/mappers/syncablemapper"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
	"github.com/figment-networks/oasishub/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(syncabledomain.Type, types.Height) bool
	Count(syncabledomain.Type) (*int64, errors.ApplicationError)
	GetByHeight(syncabledomain.Type, types.Height) (*syncabledomain.Syncable, errors.ApplicationError)
	GetMostRecent(syncabledomain.Type) (*syncabledomain.Syncable, errors.ApplicationError)
	GetMostRecentCommonHeight() (*types.Height, errors.ApplicationError)

	// Commands
	Save(*syncabledomain.Syncable) errors.ApplicationError
	Create(*syncabledomain.Syncable) errors.ApplicationError
	Upsert(syncabledomain.Type, types.Height, *syncabledomain.Syncable) errors.ApplicationError
	DeleteLast(syncabledomain.Type, int64) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) DbRepo {
	return &dbRepo{
		client: c,
	}
}

func (r *dbRepo) Exists(t syncabledomain.Type, h types.Height) bool {
	query := mainQuery(t, h)
	foundTransaction := orm.SyncableModel{}

	if err := r.client.Where(&query).Take(&foundTransaction).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count(t syncabledomain.Type) (*int64, errors.ApplicationError) {
	query := typeQuery(t)
	var count int64
	if err := r.client.Table(orm.SyncableModel{}.TableName()).Where(&query).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not get count of syncable for type %s", t), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting syncable count", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByHeight(t syncabledomain.Type, h types.Height) (*syncabledomain.Syncable, errors.ApplicationError) {
	query := mainQuery(t, h)
	foundTransaction := orm.SyncableModel{}

	if err := r.client.Where(&query).Take(&foundTransaction).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find syncable with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting syncable by height", errors.QueryError, err)
	}
	m, err := syncablemapper.FromPersistence(foundTransaction)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *dbRepo) GetMostRecent(t syncabledomain.Type) (*syncabledomain.Syncable, errors.ApplicationError) {
	q := typeQuery(t)
	foundTransaction := orm.SyncableModel{}
	if err := r.client.Where(q).Where("processed_at IS NOT NULL").Order("height desc").Take(&foundTransaction).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not find most recent syncable", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting most recent syncable", errors.QueryError, err)
	}
	m, err := syncablemapper.FromPersistence(foundTransaction)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *dbRepo) GetMostRecentCommonHeight() (*types.Height, errors.ApplicationError) {
	var syncables []*syncabledomain.Syncable
	for _, t := range syncabledomain.Types {
		s, err := r.GetMostRecent(t)
		if err != nil {
			return nil, err
		}
		syncables = append(syncables, s)
	}

	smallestH := syncables[0].Height
	for _, s := range syncables {
		if s.Height < smallestH {
			smallestH = s.Height

		}
	}

	return &smallestH, nil
}

func (r *dbRepo) Save(syncable *syncabledomain.Syncable) errors.ApplicationError {
	pr, err := syncablemapper.ToPersistence(syncable)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not save syncable", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(syncable *syncabledomain.Syncable) errors.ApplicationError {
	b, err := syncablemapper.ToPersistence(syncable)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create syncable", errors.CreateError, err)
	}
	return nil
}

func (r *dbRepo) Upsert(t syncabledomain.Type, h types.Height, syncable *syncabledomain.Syncable) errors.ApplicationError {
	query := mainQuery(t, h)
	model, err := syncablemapper.ToPersistence(syncable)
	if err != nil {
		return err
	}

	if r.Exists(t, h) {
		// Update
		if err := r.client.Where(query).Updates(model).Error; err != nil {
			log.Error(err)
			return errors.NewError("could not update syncable", errors.UpdateError, err)
		}
	} else {
		// Create
		if err := r.client.Create(model).Error; err != nil {
			log.Error(err)
			return errors.NewError("could not create syncable", errors.CreateError, err)
		}
	}
	return nil
}

func (r *dbRepo) DeleteLast(t syncabledomain.Type, offset int64) errors.ApplicationError {
	query := typeQuery(t)
	var foundTransactions []orm.SyncableModel
	var ids []int64
	if err := r.client.Where(&query).Order("id asc").Offset(offset).Find(&foundTransactions).Pluck("id", &ids).Error; err != nil {
		msg := fmt.Sprintf("could not get syncable ids with offset %d", offset)
		log.Error(err)
		return errors.NewError(msg, errors.QueryError, err)
	}

	if err := r.client.Where(ids).Delete(&orm.SyncableModel{}).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not delete syncable", errors.DeleteError, err)
	}

	return nil
}

/*************** Private ***************/

func mainQuery(t syncabledomain.Type, h types.Height) orm.SyncableModel {
	return orm.SyncableModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
		Type: t,
	}
}

func typeQuery(t syncabledomain.Type) orm.SyncableModel {
	return orm.SyncableModel{
		Type: t,
	}
}


