package getdelegationsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) ([]delegationseq.Model, errors.ApplicationError)
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

func (uc *useCase) Execute(height *types.Height) ([]delegationseq.Model, errors.ApplicationError) {
	if height == nil {
		h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
		if err != nil {
			return nil, err
		}
		height = h
	}

	ds, err := uc.delegationDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

