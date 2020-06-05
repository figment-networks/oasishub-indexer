package store

import (
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
	"time"
)

// NewIndexerMetric returns a new store from the connection string
func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	registerPlugins(conn)

	return &Store{
		db: conn,

		Database:  NewDatabaseStore(conn),
		Syncables: NewSyncablesStore(conn),
		Reports:   NewReportsStore(conn),

		BlockSeq:               NewBlockSeqStore(conn),
		DebondingDelegationSeq: NewDebondingDelegationSeqStore(conn),
		DelegationSeq:          NewDelegationSeqStore(conn),
		StakingSeq:             NewStakingSeqStore(conn),
		TransactionSeq:         NewTransactionSeqStore(conn),
		ValidatorSeq:           NewValidatorSeqStore(conn),

		AccountAgg:   NewAccountAggStore(conn),
		ValidatorAgg: NewValidatorAggStore(conn),
	}, nil
}

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Database  *DatabaseStore
	Syncables *SyncablesStore
	Reports   *ReportsStore

	BlockSeq               *BlockSeqStore
	DebondingDelegationSeq *DebondingDelegationSeqStore
	DelegationSeq          *DelegationSeqStore
	StakingSeq             *StakingSeqStore
	TransactionSeq         *TransactionSeqStore
	ValidatorSeq           *ValidatorSeqStore

	AccountAgg   *AccountAggStore
	ValidatorAgg *ValidatorAggStore
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// SetDebugMode enabled detailed query logging
func (s *Store) SetDebugMode(enabled bool) {
	s.db.LogMode(enabled)
}

// registerPlugins registers gorm plugins
func registerPlugins(c *gorm.DB) {
	c.Callback().Create().Before("gorm:create").Register("db_plugin:before_create", castQuantity)
	c.Callback().Update().Before("gorm:update").Register("db_plugin:before_update", castQuantity)
}

// castQuantity casts decimal to quantity type
func castQuantity(scope *gorm.Scope) {
	for _, f := range scope.Fields() {
		v := f.Field.Type().String()
		if v == "types.Quantity" {
			f.IsNormal = true
			t := f.Field.Interface().(types.Quantity)
			f.Field = reflect.ValueOf(gorm.Expr("cast(? AS DECIMAL(65,0))", t.String()))
		}
	}
}

func logQueryDuration(start time.Time, queryName string) {
	elapsed := time.Since(start)
	metric.DatabaseQueryDuration.WithLabelValues(queryName).Set(elapsed.Seconds())
}
