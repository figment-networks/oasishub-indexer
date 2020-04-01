package getdebondingdelegationsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/delegationdomain"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height types.Height) ([]*delegationdomain.DebondingDelegationSeq, errors.ApplicationError)
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

func (uc *useCase) Execute(height types.Height) ([]*delegationdomain.DebondingDelegationSeq, errors.ApplicationError) {
	ds, err := uc.debondingDelegationDbRepo.GetByHeight(height)
	if err != nil {
		return nil, err
	}

	return ds, nil
}


