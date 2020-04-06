package getdebondingdelegationsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/debondingdelegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*debondingdelegationseqmapper.ListView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	debondingDelegationDbRepo debondingdelegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	debondingDelegationDbRepo debondingdelegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		debondingDelegationDbRepo: debondingDelegationDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*debondingdelegationseqmapper.ListView, errors.ApplicationError) {
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

	dds, err := uc.debondingDelegationDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return debondingdelegationseqmapper.ToListView(dds), nil
}


