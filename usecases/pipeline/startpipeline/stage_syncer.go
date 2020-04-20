package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Syncer interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type syncer struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
}

func NewSyncer(syncableDbRepo syncablerepo.DbRepo, syncableProxyRepo syncablerepo.ProxyRepo) Syncer {
	return &syncer{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
	}
}

func (s *syncer) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*payload)
	h := payload.CurrentHeight

	// Sync block
	b, err := s.sync(syncable.BlockType, h)
	if err != nil {
		return nil, err
	}
	payload.BlockSyncable = b

	// Sync state
	st, err := s.sync(syncable.StateType, h)
	if err != nil {
		return nil, err
	}
	payload.StateSyncable = st

	// Sync validators
	v, err := s.sync(syncable.ValidatorsType, h)
	if err != nil {
		return nil, err
	}
	payload.ValidatorsSyncable = v

	// Sync transactions
	tr, err := s.sync(syncable.TransactionsType, h)
	if err != nil {
		return nil, err
	}
	payload.TransactionsSyncable = tr

	return payload, nil
}

/*************** Private ***************/

func (s *syncer) sync(t syncable.Type, h types.Height) (*syncable.Model, errors.ApplicationError) {
	existingSyncable, err := s.syncableDbRepo.GetByHeight(t, h)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			// Create if not found
			m, err := s.syncableProxyRepo.GetByHeight(t, h)
			if err != nil {
				return nil, err
			}

			err = s.syncableDbRepo.Create(m)
			if err != nil {
				return nil, err
			}

			return m, nil
		}
	}
	return existingSyncable, nil
}
