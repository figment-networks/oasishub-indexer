package validatorseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/mappers/validatorseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

const (
	totalSharesForIntervalQuery = `
SELECT
  time_bucket($1, time) AS time_interval,
  SUM(a) as sum,
  COUNT(*) as count,
  SUM(a) / COUNT(*) AS avg
FROM (
  SELECT
    MAX(time) as time,
    SUM(total_shares) / COUNT(*) AS a
  FROM validator_sequences
    WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $2::INTERVAL
  GROUP BY height
  ORDER BY height
) d
GROUP BY time_interval
ORDER BY time_interval;
`
	validatorSharesForIntervalQuery=`
SELECT
  time_bucket($2, time) AS time_interval,
  SUM(total_shares) / COUNT(*) AS avg
FROM validator_sequences
  WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $3::INTERVAL AND entity_uid = $1
GROUP BY time_interval
ORDER BY time_interval ASC;
`
)

var _ DbRepo = (*dbRepo)(nil)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	Count() (*int64, errors.ApplicationError)
	GetByHeight(types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError)
	GetTotalValidatedByEntityUID(types.PublicKey) (*int64, errors.ApplicationError)
	GetTotalMissedByEntityUID(types.PublicKey) (*int64, errors.ApplicationError)
	GetTotalProposedByEntityUID(types.PublicKey) (*int64, errors.ApplicationError)
	GetTotalSharesForInterval(string, string) ([]Row, errors.ApplicationError)
	GetValidatorSharesForInterval(types.PublicKey, string, string) ([]Row, errors.ApplicationError)

	// Commands
	Save(*validatordomain.ValidatorSeq) errors.ApplicationError
	Create(*validatordomain.ValidatorSeq) errors.ApplicationError
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
	query := heightQuery(h)
	foundSyncableValidator := orm.ValidatorSeqModel{}

	if err := r.client.Where(&query).Take(&foundSyncableValidator).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(orm.ValidatorSeqModel{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of validator sequences", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of validator sequences", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByHeight(h types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError) {
	query := heightQuery(h)
	var seqs []orm.ValidatorSeqModel

	if err := r.client.Where(&query).Find(&seqs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError(fmt.Sprintf("could not find validator sequences with height %d", h), errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting validator sequences", errors.QueryError, err)
	}

	var resp []*validatordomain.ValidatorSeq
	for _, s := range seqs {
		vs, err := validatorseqmapper.FromPersistence(s)
		if err != nil {
			return nil, err
		}

		resp = append(resp, vs)
	}
	return resp, nil
}

func (r *dbRepo) GetTotalValidatedByEntityUID(key types.PublicKey) (*int64, errors.ApplicationError) {
	v := true
	q := orm.ValidatorSeqModel{
		EntityUID:          key,
		PrecommitValidated: &v,
	}
	var count int64
	if err := r.client.Table(orm.ValidatorSeqModel{}.TableName()).Where(q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of validated by entityUID", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of validated by entityUID", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetTotalMissedByEntityUID(key types.PublicKey) (*int64, errors.ApplicationError) {
	v := false
	q := orm.ValidatorSeqModel{
		EntityUID:          key,
		PrecommitValidated: &v,
	}
	var count int64
	if err := r.client.Table(orm.ValidatorSeqModel{}.TableName()).Where(q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of not validated by entityUID", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of not validated by entityUID", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetTotalProposedByEntityUID(key types.PublicKey) (*int64, errors.ApplicationError) {
	q := orm.ValidatorSeqModel{
		EntityUID: key,
		Proposed:  true,
	}
	var count int64
	if err := r.client.Table(orm.ValidatorSeqModel{}.TableName()).Where(q).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("could not get count of proposed by entityUID", errors.NotFoundError, err)
		}
		log.Error(err)
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
		log.Error(err)
		return nil, errors.NewError("could not query total shares for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []Row
	for rows.Next() {
		var row Row
		if err := r.client.ScanRows(rows, &row); err != nil {
			log.Error(err)
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}

		res = append(res, row)
	}
	return res, nil
}

func (r *dbRepo) GetValidatorSharesForInterval(key types.PublicKey, interval string, period string) ([]Row, errors.ApplicationError) {
	rows, err := r.client.Raw(validatorSharesForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		log.Error(err)
		return nil, errors.NewError("could not query validator shares for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []Row
	for rows.Next() {
		var row Row
		if err := r.client.ScanRows(rows, &row); err != nil {
			log.Error(err)
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}

		res = append(res, row)
	}
	return res, nil
}

// - Commands
func (r *dbRepo) Save(sv *validatordomain.ValidatorSeq) errors.ApplicationError {
	pr, err := validatorseqmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not save validator sequence", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(sv *validatordomain.ValidatorSeq) errors.ApplicationError {
	b, err := validatorseqmapper.ToPersistence(sv)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create validator sequence", errors.CreateError, err)
	}
	return nil
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.ValidatorSeqModel {
	return orm.ValidatorSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}
