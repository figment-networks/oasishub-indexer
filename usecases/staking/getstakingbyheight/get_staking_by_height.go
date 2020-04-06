package getstakingbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/stakingseqmapper"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*stakingseqmapper.DetailsView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	stakingDbRepo stakingseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	stakingDbRepo stakingseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		stakingDbRepo: stakingDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*stakingseqmapper.DetailsView, errors.ApplicationError) {
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

	ss, err := uc.stakingDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return stakingseqmapper.ToDetailsView(ss), nil
}

