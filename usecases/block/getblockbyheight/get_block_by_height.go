package getblockbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height types.Height) (*syncabledomain.Syncable, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo   syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	blockDbRepo      blockseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	blockDbRepo blockseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:   syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		blockDbRepo:      blockDbRepo,
	}
}

func (uc *useCase) Execute(height types.Height) (*syncabledomain.Syncable, errors.ApplicationError) {
	// First check for syncable in database then if not found get it from node
	syncable, err := uc.syncableDbRepo.GetByHeight(syncabledomain.BlockType, height)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			syncable, err = uc.syncableProxyRepo.GetByHeight(syncabledomain.BlockType, height)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return syncable, nil
}
