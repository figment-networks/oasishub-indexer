package debondingdelegationseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	GetByHeight(types.Height) ([]debondingdelegationseq.Model, errors.ApplicationError)
	GetRecentByValidatorUID(types.PublicKey, int64) ([]debondingdelegationseq.Model, errors.ApplicationError)
	GetRecentByDelegatorUID(types.PublicKey, int64) ([]debondingdelegationseq.Model, errors.ApplicationError)

	// Commands
	Create(*debondingdelegationseq.Model) errors.ApplicationError
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
	m := debondingdelegationseq.Model{}

	if err := r.client.Where(&q).Find(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) GetByHeight(h types.Height) ([]debondingdelegationseq.Model, errors.ApplicationError) {
	q := heightQuery(h)
	var ms []debondingdelegationseq.Model

	if err := r.client.Where(&q).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find debinding delegation with height %d", h), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting debonding delegation by height", errors.QueryError, err)
	}
	return ms, nil
}

func (r *dbRepo) GetRecentByValidatorUID(key types.PublicKey, limit int64) ([]debondingdelegationseq.Model, errors.ApplicationError) {
	q := debondingdelegationseq.Model{
		ValidatorUID:  key,
	}
	var ms []debondingdelegationseq.Model

	if err := r.client.Where(&q).Order("height DESC").Limit(limit).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find debinding delegation with for validator %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting debonding delegation for validator", errors.QueryError, err)
	}
	return ms, nil
}

func (r *dbRepo) GetRecentByDelegatorUID(key types.PublicKey, limit int64) ([]debondingdelegationseq.Model, errors.ApplicationError) {
	q := debondingdelegationseq.Model{
		DelegatorUID:  key,
	}
	var ms []debondingdelegationseq.Model

	if err := r.client.Where(&q).Order("height DESC").Limit(limit).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find debinding delegation with for validator %s", key), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting debonding delegation for validator", errors.QueryError, err)
	}
	return ms, nil
}

// - Commands
func (r *dbRepo) Create(m *debondingdelegationseq.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		return errors.NewError("could not create debonding delegation", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) debondingdelegationseq.Model {
	return debondingdelegationseq.Model{
		Sequence: &shared.Sequence{
			Height: h,
		},
	}
}
