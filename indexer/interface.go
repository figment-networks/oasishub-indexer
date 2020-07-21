package indexer

import "github.com/figment-networks/oasishub-indexer/model"

type AccountAggStore interface {
	FindByPublicKey(key string) (*model.AccountAgg, error)
	Create(record interface{}) error
	Save(record interface{}) error
}

type SyncableStore interface {
	FindMostRecent() (*model.Syncable, error)
	CreateOrUpdate(val *model.Syncable) error
}

type ValidatorAggStore interface {
	FindByEntityUID(key string) (*model.ValidatorAgg, error)
	Create(record interface{}) error
	Save(record interface{}) error
}
