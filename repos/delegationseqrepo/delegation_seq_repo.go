package delegationseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/delegationdomain"
	"github.com/figment-networks/oasishub-indexer/mappers/delegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	GetByHeight(types.Height) ([]*delegationdomain.DelegationSeq, errors.ApplicationError)
	GetLastByValidatorUID(types.PublicKey) ([]*delegationdomain.DelegationSeq, errors.ApplicationError)
	GetLastByDelegatorUID(types.PublicKey) ([]*delegationdomain.DelegationSeq, errors.ApplicationError)

	// Commands
	Create(*delegationdomain.DelegationSeq) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) DbRepo {
	return &dbRepo{
		client: c,
	}
}

// - Queries
func (r *dbRepo) Exists(h types.Height) bool {
	query := heightQuery(h)
	foundBlock := orm.DelegationSeqModel{}

	if err := r.client.Where(&query).Take(&foundBlock).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) GetByHeight(h types.Height) ([]*delegationdomain.DelegationSeq, errors.ApplicationError) {
	query := heightQuery(h)
	var seqs []orm.DelegationSeqModel

	if err := r.client.Where(&query).Find(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find delegation sequence with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting delegation sequence by height", errors.QueryError, err)
	}

	var resp []*delegationdomain.DelegationSeq
	for _, s := range seqs {
		vs, err := delegationseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, vs)
	}
	return resp, nil
}

func (r *dbRepo) GetLastByValidatorUID(key types.PublicKey) ([]*delegationdomain.DelegationSeq, errors.ApplicationError) {
	query := orm.DelegationSeqModel{
		ValidatorUID:  key,
	}
	var seqs []orm.DelegationSeqModel

	sub := r.client.Table(orm.DelegationSeqModel{}.TableName()).Select("height").Order("height DESC").Limit(1).QueryExpr()
	if err := r.client.Where(&query).Where("height = (?)", sub).Find(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find delegation sequence for key %s", key), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting delegation sequence by key", errors.QueryError, err)
	}

	var resp []*delegationdomain.DelegationSeq
	for _, s := range seqs {
		vs, err := delegationseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, vs)
	}
	return resp, nil
}

func (r *dbRepo) GetLastByDelegatorUID(key types.PublicKey) ([]*delegationdomain.DelegationSeq, errors.ApplicationError) {
	query := orm.DelegationSeqModel{
		DelegatorUID:  key,
	}
	var seqs []orm.DelegationSeqModel

	sub := r.client.Table(orm.DelegationSeqModel{}.TableName()).Select("height").Order("height DESC").Limit(1).QueryExpr()
	if err := r.client.Where(&query).Where("height = (?)", sub).Find(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find delegation sequence for key %s", key), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting delegation sequence by key", errors.QueryError, err)
	}

	var resp []*delegationdomain.DelegationSeq
	for _, s := range seqs {
		vs, err := delegationseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, vs)
	}
	return resp, nil
}

// - Mutations
func (r *dbRepo) Create(block *delegationdomain.DelegationSeq) errors.ApplicationError {
	b, err := delegationseqmapper.ToPersistence(block)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create delegation sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.DelegationSeqModel {
	return orm.DelegationSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}
