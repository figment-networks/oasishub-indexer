package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/repos/syncablerepo"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
	"github.com/figment-networks/oasishub/utils/pipeline"
)

type Syncer interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type syncer struct {
	syncableDbRepo      syncablerepo.DbRepo
	syncableProxyRepo    syncablerepo.ProxyRepo
}

func NewSyncer(syncableDbRepo syncablerepo.DbRepo, syncableProxyRepo syncablerepo.ProxyRepo) Syncer {
	return &syncer{
		syncableDbRepo:      syncableDbRepo,
		syncableProxyRepo:    syncableProxyRepo,
	}
}

func (s *syncer) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*payload)
	h := payload.CurrentHeight

	// Sync block
	b, err := s.sync(syncabledomain.BlockType, h)
	if err != nil {
		return nil, err
	}
	payload.BlockSyncable = b

	// Sync state
	st, err := s.sync(syncabledomain.StateType, h)
	if err != nil {
		return nil, err
	}
	payload.StateSyncable = st

	// Sync validators
	v, err := s.sync(syncabledomain.ValidatorsType, h)
	if err != nil {
		return nil, err
	}
	payload.ValidatorsSyncable = v

	// Sync transactions
	tr, err := s.sync(syncabledomain.TransactionsType, h)
	if err != nil {
		return nil, err
	}
	payload.TransactionsSyncable = tr

	return payload, nil
}

/*************** Private ***************/

func (s *syncer) sync(t syncabledomain.Type, h types.Height) (*syncabledomain.Syncable, errors.ApplicationError) {
	existingSyncable, err := s.syncableDbRepo.GetByHeight(t, h)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			// Create if not found
			syncable, err := s.syncableProxyRepo.GetByHeight(t, h)
			if err != nil {
				return nil, err
			}

			err = s.syncableDbRepo.Create(syncable)
			if err != nil {
				return nil, err
			}

			return syncable, nil
		}
	}

	return existingSyncable, nil
}
