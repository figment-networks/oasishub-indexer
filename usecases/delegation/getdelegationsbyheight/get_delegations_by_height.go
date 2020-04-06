package getdelegationsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/delegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*delegationseqmapper.ListView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	delegationDbRepo delegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	delegationDbRepo delegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		delegationDbRepo: delegationDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*delegationseqmapper.ListView, errors.ApplicationError) {
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

	ds, err := uc.delegationDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return delegationseqmapper.ToListView(ds), nil
}

