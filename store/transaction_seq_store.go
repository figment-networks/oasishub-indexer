package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

var (
	_ TransactionSeqStore = (*transactionSeqStore)(nil)
)

type TransactionSeqStore interface {
	BaseStore

	FindByHeight(h int64) ([]model.TransactionSeq, error)
}

func NewTransactionSeqStore(db *gorm.DB) *transactionSeqStore {
	return &transactionSeqStore{scoped(db, model.TransactionSeq{})}
}

// transactionSeqStore handles operations on transactions
type transactionSeqStore struct {
	baseStore
}

func (s transactionSeqStore) FindByHeight(h int64) ([]model.TransactionSeq, error) {
	q := model.TransactionSeq{
		Sequence: &model.Sequence{
			Height: h,
		},
	}
	var result []model.TransactionSeq

	err := s.db.
		Where(&q).
		Find(&result).
		Error

	return result, checkErr(err)
}


