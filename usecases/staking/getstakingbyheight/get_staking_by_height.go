package getstakingbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/stakingdomain"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*stakingdomain.StakingSeq, errors.ApplicationError)
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

func (uc *useCase) Execute(height *types.Height) (*stakingdomain.StakingSeq, errors.ApplicationError) {
	if height == nil {
		h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
		if err != nil {
			return nil, err
		}
		height = h
	}

	es, err := uc.stakingDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return es, nil
}

