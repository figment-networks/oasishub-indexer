package debondingdelegationseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/delegationdomain"
	"github.com/figment-networks/oasishub-indexer/mappers/debondingdelegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	GetByHeight(types.Height) ([]*delegationdomain.DebondingDelegationSeq, errors.ApplicationError)

	// Commands
	Create(*delegationdomain.DebondingDelegationSeq) errors.ApplicationError
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
	foundBlock := orm.DebondingDelegationSeqModel{}

	if err := r.client.Where(&query).Take(&foundBlock).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) GetByHeight(h types.Height) ([]*delegationdomain.DebondingDelegationSeq, errors.ApplicationError) {
	query := heightQuery(h)
	var seqs []orm.DebondingDelegationSeqModel

	if err := r.client.Where(&query).Find(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find debinding delegation with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting debonding delegation by height", errors.QueryError, err)
	}

	var resp []*delegationdomain.DebondingDelegationSeq
	for _, s := range seqs {
		vs, err := debondingdelegationseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, vs)
	}
	return resp, nil
}

func (r *dbRepo) Create(block *delegationdomain.DebondingDelegationSeq) errors.ApplicationError {
	b, err := debondingdelegationseqmapper.ToPersistence(block)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create debonding delegation", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.DebondingDelegationSeqModel {
	return orm.DebondingDelegationSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}
