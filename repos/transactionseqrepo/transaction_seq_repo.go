package transactionseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/transactiondomain"
	"github.com/figment-networks/oasishub-indexer/mappers/transactionseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	Count() (*int64, errors.ApplicationError)
	GetByHeight(types.Height) ([]*transactiondomain.TransactionSeq, errors.ApplicationError)

	// Commands
	Save(*transactiondomain.TransactionSeq) errors.ApplicationError
	Create(*transactiondomain.TransactionSeq) errors.ApplicationError
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
	foundTransaction := orm.TransactionSeqModel{}

	if err := r.client.Where(&query).Take(&foundTransaction).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(orm.TransactionSeqModel{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not get count of transaction sequences"), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of transaction sequences", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByHeight(h types.Height) ([]*transactiondomain.TransactionSeq, errors.ApplicationError) {
	query := heightQuery(h)
	var seqs []orm.TransactionSeqModel

	if err := r.client.Where(&query).Take(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find transaction sequences with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting transaction sequences by height", errors.QueryError, err)
	}

	var resp []*transactiondomain.TransactionSeq
	for _, s := range seqs {
		ts, err := transactionseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, ts)
	}
	return resp, nil
}

func (r *dbRepo) Save(transaction *transactiondomain.TransactionSeq) errors.ApplicationError {
	pr, err := transactionseqmapper.ToPersistence(transaction)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not save transaction sequence", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(transaction *transactiondomain.TransactionSeq) errors.ApplicationError {
	b, err := transactionseqmapper.ToPersistence(transaction)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create transaction sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.TransactionSeqModel {
	return orm.TransactionSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}


