package getblockbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/blockseqmapper"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*blockseqmapper.DetailsView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo       syncablerepo.DbRepo
	syncableProxyRepo    syncablerepo.ProxyRepo
	blockSeqDbRepo       blockseqrepo.DbRepo
	validatorSeqDbRepo   validatorseqrepo.DbRepo
	transactionSeqDbRepo transactionseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	blockSeqDbRepo blockseqrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	transactionSeqDbRepo transactionseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:       syncableDbRepo,
		syncableProxyRepo:    syncableProxyRepo,
		blockSeqDbRepo:       blockSeqDbRepo,
		validatorSeqDbRepo:   validatorSeqDbRepo,
		transactionSeqDbRepo: transactionSeqDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*blockseqmapper.DetailsView, errors.ApplicationError) {
	// Get last indexed height
	lastH, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
	if err != nil {
		return nil, err
	}

	// Show last synced height, if not provided
	if height == nil {
		height = lastH
	}

	if height.Larger(*lastH) {
		return nil, errors.NewErrorFromMessage("height is not indexed", errors.ServerInvalidParamsError)
	}

	// First check if syncable is in database
	s, err := uc.syncableDbRepo.GetByHeight(syncable.BlockType, *height)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			// If it is not there get it from proxy
			s, err = uc.syncableProxyRepo.GetByHeight(syncable.BlockType, *height)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	bs, err := uc.blockSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	vs, err := uc.validatorSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	ts, err := uc.transactionSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return blockseqmapper.ToDetailsView(bs, *s, vs, ts)
}
