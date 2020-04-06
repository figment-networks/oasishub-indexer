package delegationseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	GetByHeight(types.Height) ([]delegationseq.Model, errors.ApplicationError)
	GetLastByValidatorUID(types.PublicKey) ([]delegationseq.Model, errors.ApplicationError)
	GetCurrentByDelegatorUID(types.PublicKey) ([]delegationseq.Model, errors.ApplicationError)

	// Commands
	Create(*delegationseq.Model) errors.ApplicationError
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
	q := heightQuery(h)
	m := delegationseq.Model{}

	if err := r.client.Where(&q).First(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) GetByHeight(h types.Height) ([]delegationseq.Model, errors.ApplicationError) {
	q := heightQuery(h)
	var ms []delegationseq.Model

	if err := r.client.Where(&q).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find delegation sequence with height %d", h), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting delegation sequence by height", errors.QueryError, err)
	}
	return ms, nil
}

func (r *dbRepo) GetLastByValidatorUID(key types.PublicKey) ([]delegationseq.Model, errors.ApplicationError) {
	q := delegationseq.Model{
		ValidatorUID:  key,
	}
	var ms []delegationseq.Model

	sub := r.client.Table(delegationseq.Model{}.TableName()).Select("height").Order("height DESC").Limit(1).QueryExpr()
	if err := r.client.Where(&q).Where("height = (?)", sub).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find delegation sequence for key %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting delegation sequence by key", errors.QueryError, err)
	}
	return ms, nil
}

func (r *dbRepo) GetCurrentByDelegatorUID(key types.PublicKey) ([]delegationseq.Model, errors.ApplicationError) {
	q := delegationseq.Model{
		DelegatorUID:  key,
	}
	var ms []delegationseq.Model

	sub := r.client.Table(delegationseq.Model{}.TableName()).Select("height").Order("height DESC").Limit(1).QueryExpr()
	if err := r.client.Where(&q).Where("height = (?)", sub).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find delegation sequence for key %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting delegation sequence by key", errors.QueryError, err)
	}
	return ms, nil
}

// - Mutations
func (r *dbRepo) Create(m *delegationseq.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create delegation sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) delegationseq.Model {
	return delegationseq.Model{
		Sequence: &shared.Sequence{
			Height: h,
		},
	}
}
