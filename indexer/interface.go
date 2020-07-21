package indexer

import "github.com/figment-networks/oasishub-indexer/model"

type AccountAggStore interface {
	FindByPublicKey(key string) (*model.AccountAgg, error)
	Create(record interface{}) error
	Save(record interface{}) error
}

type BlockSeqStore interface {
	CreateIfNotExists(block *model.BlockSeq) error
}

type DebondingDelegationSeqStore interface {
	FindByHeight(h int64) ([]model.DebondingDelegationSeq, error)
	Create(record interface{}) error
}

type DelegationSeqStore interface {
	FindByHeight(h int64) ([]model.DelegationSeq, error)
	Create(record interface{}) error
}

type StakingSeqStore interface {
	FindByHeight(height int64) (*model.StakingSeq, error)
	Create(record interface{}) error
}

type SyncableStore interface {
	FindMostRecent() (*model.Syncable, error)
	CreateOrUpdate(val *model.Syncable) error
}

type TransactionSeqStore interface {
	FindByHeight(h int64) ([]model.TransactionSeq, error)
	Create(record interface{}) error
}

type ValidatorAggStore interface {
	FindByEntityUID(key string) (*model.ValidatorAgg, error)
	Create(record interface{}) error
	Save(record interface{}) error
}

type ValidatorSeqStore interface {
	FindByHeight(h int64) ([]model.ValidatorSeq, error)
	Create(record interface{}) error
}
