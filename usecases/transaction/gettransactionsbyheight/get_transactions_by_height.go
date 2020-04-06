package gettransactionsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/transactionseqmapper"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*transactionseqmapper.ListView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	transactionDbRepo transactionseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	transactionDbRepo transactionseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		transactionDbRepo: transactionDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*transactionseqmapper.ListView, errors.ApplicationError) {
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

	ts, err := uc.transactionDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return transactionseqmapper.ToListView(ts), nil
}
