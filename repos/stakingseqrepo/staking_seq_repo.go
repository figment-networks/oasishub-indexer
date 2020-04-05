package stakingseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/stakingseq"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	GetByHeight(types.Height) (*stakingseq.Model, errors.ApplicationError)

	// Commands
	Create(*stakingseq.Model) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) DbRepo {
	return &dbRepo{
		client: c,
	}
}

func (r *dbRepo) Exists(h types.Height) bool {
	q := heightQuery(h)
	m := stakingseq.Model{}

	if err := r.client.Where(&q).Find(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) GetByHeight(h types.Height) (*stakingseq.Model, errors.ApplicationError) {
	q := heightQuery(h)
	var m stakingseq.Model

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find staking sequence with height %d", h), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting staking sequence by height", errors.QueryError, err)
	}
	return &m, nil
}

func (r *dbRepo) Create(m *stakingseq.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		return errors.NewError("could not create staking sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) stakingseq.Model {
	return stakingseq.Model{
		Sequence: &shared.Sequence{
			Height: h,
		},
	}
}
