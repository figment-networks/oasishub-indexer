package blockseqrepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/blockdomain"
	"github.com/figment-networks/oasishub-indexer/mappers/blockseqmapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
)

const (
	blockTimesForRecentBlocksQuery = `
SELECT 
  MIN(height) start_height, 
  MAX(height) end_height, 
  MIN(time) start_time,
  MAX(time) end_time,
  COUNT(*) count, 
  EXTRACT(EPOCH FROM MAX(time) - MIN(time)) AS diff, 
  EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg
  FROM ( 
    SELECT * FROM block_sequences
    ORDER BY height DESC
    LIMIT ?
  ) t;
`
	blockTimesForIntervalQuery = `
SELECT 
  time_bucket(?, time) AS interval, 
  COUNT(*) AS count, 
  EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg
FROM block_sequences
GROUP BY interval;
`
)

type DbRepo interface {
	// Queries
	Exists(types.Height) bool
	Count() (*int64, errors.ApplicationError)
	GetByHeight(types.Height) (*blockdomain.BlockSeq, errors.ApplicationError)
	GetMostRecent(BlockDbQuery) (*blockdomain.BlockSeq, errors.ApplicationError)
	GetAvgBlockTimesForRecentBlocks(int64) Result
	GetAvgBlockTimesForInterval(string) ([]Row, errors.ApplicationError)

	// Commands
	Save(*blockdomain.BlockSeq) errors.ApplicationError
	Create(*blockdomain.BlockSeq) errors.ApplicationError
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
	foundBlock := orm.BlockSeqModel{}

	if err := r.client.Where(&query).Take(&foundBlock).Error; err != nil {
		return false
	}
	return true
}

func (r *dbRepo) Count() (*int64, errors.ApplicationError) {
	var count int64
	if err := r.client.Table(orm.BlockSeqModel{}.TableName()).Count(&count).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("block sequence not found", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("error getting count of block sequences", errors.QueryError, err)
	}

	return &count, nil
}

func (r *dbRepo) GetByHeight(h types.Height) (*blockdomain.BlockSeq, errors.ApplicationError) {
	query := heightQuery(h)
	seq := orm.BlockSeqModel{}

	if err := r.client.Where(&query).Take(&seq).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("block sequence not found", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError(fmt.Sprintf("could not find block sequence with height %d", h), errors.QueryError, err)
	}
	m, err := blockseqmapper.FromPersistence(seq)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *dbRepo) GetMostRecent(q BlockDbQuery) (*blockdomain.BlockSeq, errors.ApplicationError) {
	seq := orm.BlockSeqModel{}
	if err := r.client.Where(q.String()).Order("height desc").Take(&seq).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewError("most recent block sequence not found", errors.NotFoundError, err)
		}
		log.Error(err)
		return nil, errors.NewError("could not find most recent block sequence", errors.QueryError, err)
	}
	m, err := blockseqmapper.FromPersistence(seq)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type Result struct {
	StartHeight int64   `json:"start_height"`
	EndHeight   int64   `json:"end_height"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Count       int64   `json:"count"`
	Diff        float64 `json:"diff"`
	Avg         float64 `json:"avg"`
}

func (r *dbRepo) GetAvgBlockTimesForRecentBlocks(limit int64) Result {
	var result Result
	r.client.Raw(blockTimesForRecentBlocksQuery, limit).Scan(&result)

	return result
}

type Row struct {
	Interval string  `json:"interval"`
	Count    int64   `json:"count"`
	Avg      float64 `json:"avg"`
}

func (r *dbRepo) GetAvgBlockTimesForInterval(interval string) ([]Row, errors.ApplicationError) {
	rows, err := r.client.Raw(blockTimesForIntervalQuery, interval).Rows()
	if err != nil {
		log.Error(err)
		return nil, errors.NewError("could not query block times for interval", errors.QueryError, err)
	}
	defer rows.Close()

	var res []Row
	for rows.Next() {
		var row Row
		if err := r.client.ScanRows(rows, &row); err != nil {
			log.Error(err, log.Field("resource", "block_seq"))
			return nil, errors.NewError("could not scan rows", errors.QueryError, err)
		}

		res = append(res, row)
	}
	return res, nil
}

func (r *dbRepo) Save(block *blockdomain.BlockSeq) errors.ApplicationError {
	pr, err := blockseqmapper.ToPersistence(block)
	if err != nil {
		return err
	}

	if err := r.client.Save(pr).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not save block sequence", errors.SaveError, err)
	}
	return nil
}

func (r *dbRepo) Create(block *blockdomain.BlockSeq) errors.ApplicationError {
	b, err := blockseqmapper.ToPersistence(block)
	if err != nil {
		return err
	}

	if err := r.client.Create(b).Error; err != nil {
		log.Error(err)
		return errors.NewError("could not create block sequence", errors.CreateError, err)
	}
	return nil
}

type BlockDbQuery struct {
	Processed bool
}

func (q *BlockDbQuery) String() string {
	query := ""
	if q.Processed {
		query += "processed_at IS NOT NULL"
	}
	return query
}

/*************** Private ***************/

func heightQuery(h types.Height) orm.BlockSeqModel {
	return orm.BlockSeqModel{
		SequenceModel: orm.SequenceModel{
			Height: h,
		},
	}
}