package stakingseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/stakingdomain"
	"github.com/figment-networks/oasishub-indexer/mappers/stakingseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	GetByHeight(types.Height) (*stakingdomain.StakingSeq, errors.ApplicationError)

	// Commands
	Create(*stakingdomain.StakingSeq) errors.ApplicationError
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
	query := heightQuery(h)
	foundBlock := orm.StakingSeqModel{}

	if err := r.client.Where(&query).Take(&foundBlock).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) GetByHeight(h types.Height) (*stakingdomain.StakingSeq, errors.ApplicationError) {
	query := heightQuery(h)
	var seq orm.StakingSeqModel

	if err := r.client.Where(&query).Take(&seq).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find staking sequence with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting staking sequence by height", errors.QueryError, err)
	}
	m, err := stakingseqmapper.FromPersistence(seq)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *dbRepo) Create(block *stakingdomain.StakingSeq) errors.ApplicationError {
	b, err := stakingseqmapper.ToPersistence(block)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create staking sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.StakingSeqModel {
	return orm.StakingSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}
