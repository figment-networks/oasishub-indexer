package validatorseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/validatorseq"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/jinzhu/gorm"
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	Count() (*int64, errors.ApplicationError)
	GetByHeight(types.Height) ([]validatorseq.Model, errors.ApplicationError)
	GetTotalValidatedByEntityUID(types.PublicKey) (*int64, errors.ApplicationError)
	GetTotalMissedByEntityUID(types.PublicKey) (*int64, errors.ApplicationError)
	GetTotalProposedByEntityUID(types.PublicKey) (*int64, errors.ApplicationError)
	GetTotalSharesForInterval(string, string) ([]Row, errors.ApplicationError)
	GetTotalVotingPowerForInterval(string, string) ([]FloatRow, errors.ApplicationError)
	GetValidatorSharesForInterval(types.PublicKey, string, string) ([]Row, errors.ApplicationError)
	GetValidatorVotingPowerForInterval(types.PublicKey, string, string) ([]FloatRow, errors.ApplicationError)
	GetValidatorUptimeForInterval(types.PublicKey, string, string) ([]FloatRow, errors.ApplicationError)

	// Commands
	Save(*validatorseq.Model) errors.ApplicationError
	Create(*validatorseq.Model) errors.ApplicationError
}

type dbRepo struct {
	client *gorm.DB
}

func NewDbRepo(c *gorm.DB) *dbRepo {
	return &dbRepo{
		client: c,
	}
}

// - Queries
func (r *dbRepo) Exists(h types.Height) bool {
	q := heightQuery(h)
	m := validatorseq.Model{}

	if err := r.client.Where(&q).Find(&m).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(validatorseq.Model{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of validator sequences", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of validator sequences", errors.QueryError, err)
	}
	return &count, nil
}

func (r *dbRepo) GetByHeight(h types.Height) ([]validatorseq.Model, errors.ApplicationError) {
	q := heightQuery(h)
	var ms []validatorseq.Model

	if err := r.client.Where(&q).Find(&ms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find validator sequences with height %d", h), errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting validator sequences", errors.QueryError, err)
	}
	return ms, nil
}

func (r *dbRepo) GetTotalValidatedByEntityUID(key types.PublicKey) (*int64, errors.ApplicationError) {
	v := true
	q := validatorseq.Model{
		EntityUID:          key,
		PrecommitValidated: &v,
	}
	var count int64
	if err := r.client.Table(validatorseq.Model{}.TableName()).Where(&q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of validated by entityUID", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of validated by entityUID", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetTotalMissedByEntityUID(key types.PublicKey) (*int64, errors.ApplicationError) {
	v := false
	q := validatorseq.Model{
		EntityUID:          key,
		PrecommitValidated: &v,
	}
	var count int64
	if err := r.client.Table(validatorseq.Model{}.TableName()).Where(&q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of not validated by entityUID", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of not validated by entityUID", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetTotalProposedByEntityUID(key types.PublicKey) (*int64, errors.ApplicationError) {
	q := validatorseq.Model{
		EntityUID: key,
		Proposed:  true,
	}
	var count int64
	if err := r.client.Table(validatorseq.Model{}.TableName()).Where(&q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of proposed by entityUID", errors.NotFoundError, err)
		}
		return nil, errors.NewError("error getting count of proposed by entityUID", errors.QueryError, err)
	}

	return &count, nil
}

type Row struct {
	TimeInterval string         `json:"time_interval"`
	Avg          types.Quantity `json:"avg"`
}

func (r *dbRepo) GetTotalSharesForInterval(interval string, period string) ([]Row, errors.ApplicationError) {
	rows, err := r.client.Raw(totalSharesForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, errors.NewError("could not query total shares for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []Row
	for rows.Next() {
		var row Row
		if err := r.client.ScanRows(rows, &row); err != nil {
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (r *dbRepo) GetTotalVotingPowerForInterval(interval string, period string) ([]FloatRow, errors.ApplicationError) {
	rows, err := r.client.Raw(totalVotingPowerForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, errors.NewError("could not query total voting power for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []FloatRow
	for rows.Next() {
		var row FloatRow
		if err := r.client.ScanRows(rows, &row); err != nil {
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}
		res = append(res, row)
	}
	return res, nil
}

type FloatRow struct {
	TimeInterval string  `json:"time_interval"`
	Avg          float64 `json:"avg"`
}

func (r *dbRepo) GetValidatorSharesForInterval(key types.PublicKey, interval string, period string) ([]Row, errors.ApplicationError) {
	rows, err := r.client.Raw(validatorSharesForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		return nil, errors.NewError("could not query validator shares for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []Row
	for rows.Next() {
		var row Row
		if err := r.client.ScanRows(rows, &row); err != nil {
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (r *dbRepo) GetValidatorVotingPowerForInterval(key types.PublicKey, interval string, period string) ([]FloatRow, errors.ApplicationError) {
	rows, err := r.client.Debug().Raw(validatorVotingPowerForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		return nil, errors.NewError("could not query validator voting power for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []FloatRow
	for rows.Next() {
		var row FloatRow
		if err := r.client.ScanRows(rows, &row); err != nil {
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (r *dbRepo) GetValidatorUptimeForInterval(key types.PublicKey, interval string, period string) ([]FloatRow, errors.ApplicationError) {
	rows, err := r.client.Raw(validatorUptimeForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		return nil, errors.NewError("could not query validator uptime for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []FloatRow
	for rows.Next() {
		var row FloatRow
		if err := r.client.ScanRows(rows, &row); err != nil {
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}
		res = append(res, row)
	}
	return res, nil
}

// - Commands
func (r *dbRepo) Save(m *validatorseq.Model) errors.ApplicationError {
	if err := r.client.Save(m).Error; err != nil {
		return errors.NewError("could not save validator sequence", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(m *validatorseq.Model) errors.ApplicationError {
	if err := r.client.Create(m).Error; err != nil {
		return errors.NewError("could not create validator sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) validatorseq.Model {
	return validatorseq.Model{
		Sequence: &shared.Sequence{
			Height: h,
		},
	}
}
