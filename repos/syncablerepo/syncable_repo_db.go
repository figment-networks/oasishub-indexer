package syncablerepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(syncable.Type, types.Height) bool
	Count(syncable.Type) (*int64, errors.ApplicationError)
	GetByHeight(syncable.Type, types.Height) (*syncable.Model, errors.ApplicationError)
	GetMostRecent(syncable.Type) (*syncable.Model, errors.ApplicationError)
	GetMostRecentCommonHeight() (*types.Height, errors.ApplicationError)

	// Commands
	Save(*syncable.Model) errors.ApplicationError
	Create(*syncable.Model) errors.ApplicationError
	Upsert(syncable.Type, types.Height, *syncable.Model) errors.ApplicationError
	DeletePrevByHeight(types.Height) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) DbRepo {
	return &dbRepo{
		client: c,
	}
}

func (r *dbRepo) Exists(t syncable.Type, h types.Height) bool {
	query := mainQuery(t, h)
	foundTransaction := syncable.Model{}

	if err := r.client.Where(&query).First(&foundTransaction).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count(t syncable.Type) (*int64, errors.ApplicationError) {
	q := typeQuery(t)
	var count int64
	if err := r.client.Table(syncable.Model{}.TableName()).Where(&q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not get count of syncable for type %s", t), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting syncable count", errors.QueryError, err)
	}
	return &count, nil
}

func (r *dbRepo) GetByHeight(t syncable.Type, h types.Height) (*syncable.Model, errors.ApplicationError) {
	q := mainQuery(t, h)
	m := syncable.Model{}

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find syncable with height %d", h), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting syncable by height", errors.QueryError, err)
	}
	return &m, nil
}

func (r *dbRepo) GetMostRecent(t syncable.Type) (*syncable.Model, errors.ApplicationError) {
	q := typeQuery(t)
	m := syncable.Model{}
	if err := r.client.Where(&q).Where("processed_at IS NOT NULL").Order("height desc").First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not find most recent syncable", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting most recent syncable", errors.QueryError, err)
	}
	return &m, nil
}

func (r *dbRepo) GetMostRecentCommonHeight() (*types.Height, errors.ApplicationError) {
	var syncables []*syncable.Model
	for _, t := range syncable.Types {
		s, err := r.GetMostRecent(t)
		if err != nil {
			return nil, err
		}
		syncables = append(syncables, s)
	}

	// If there are not syncables yet, just start from the beginning
	if len(syncables) == 0 {
		h := config.FirstBlockHeight()
		return &h, nil
	}

	smallestH := syncables[0].Height
	for _, s := range syncables {
		if s.Height < smallestH {
			smallestH = s.Height

		}
	}
	return &smallestH, nil
}

func (r *dbRepo) Save(m *syncable.Model) errors.ApplicationError {
	if err := r.client.Save(m).Error; err != nil {
		return errors.NewError("could not save syncable", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(m *syncable.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		return errors.NewError("could not create syncable", errors.CreateError, err)
	}
	return nil
}

func (r *dbRepo) Upsert(t syncable.Type, h types.Height, m *syncable.Model) errors.ApplicationError {
	q := mainQuery(t, h)

	if r.Exists(t, h) {
		// Update
		if err := r.client.Where(&q).Updates(m).Error; err != nil {
			return errors.NewError("could not update syncable", errors.UpdateError, err)
		}
	} else {
		// Create
		if err := r.client.Create(m).Error; err != nil {
			return errors.NewError("could not create syncable", errors.CreateError, err)
		}
	}
	return nil
}

func (r *dbRepo) DeletePrevByHeight(maxHeight types.Height) errors.ApplicationError {
	if err := r.client.Debug().Where("height <= ?", maxHeight).Delete(&syncable.Model{}).Error; err != nil {
		return errors.NewError("could not delete syncables", errors.DeleteError, err)
	}
	return nil
}

/*************** Private ***************/

func mainQuery(t syncable.Type, h types.Height) syncable.Model {
	return syncable.Model{
		Sequence: &shared.Sequence{
			Height: h,
		},
		Type: t,
	}
}

func typeQuery(t syncable.Type) syncable.Model {
	return syncable.Model{
		Type: t,
	}
}


