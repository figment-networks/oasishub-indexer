package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

func NewTransactionSeqStore(db *gorm.DB) *TransactionSeqStore {
	return &TransactionSeqStore{scoped(db, model.TransactionSeq{})}
}

// TransactionSeqStore handles operations on transactions
type TransactionSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the transaction if it does not exist
func (s TransactionSeqStore) CreateIfNotExists(transaction *model.TransactionSeq) error {
	_, err := s.FindByHeight(transaction.Height)
	if isNotFound(err) {
		return s.Create(transaction)
	}
	return nil
}

func (s TransactionSeqStore) FindByHeight(h int64) ([]model.TransactionSeq, error) {
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


