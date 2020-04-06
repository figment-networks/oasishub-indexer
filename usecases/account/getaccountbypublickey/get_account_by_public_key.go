package getaccountbypublickey

import (
	"github.com/figment-networks/oasishub-indexer/mappers/accountaggmapper"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(key types.PublicKey) (*accountaggmapper.DetailsView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo               syncablerepo.DbRepo
	syncableProxyRepo            syncablerepo.ProxyRepo
	accountSeqDbRepo             accountaggrepo.DbRepo
	delegationSeqDbRepo          delegationseqrepo.DbRepo
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	accountSeqDbRepo accountaggrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		accountSeqDbRepo:  accountSeqDbRepo,
		delegationSeqDbRepo: delegationSeqDbRepo,
		debondingDelegationSeqDbRepo: debondingDelegationSeqDbRepo,
	}
}

func (uc *useCase) Execute(key types.PublicKey) (*accountaggmapper.DetailsView, errors.ApplicationError) {
	aa, err := uc.accountSeqDbRepo.GetByPublicKey(key)
	if err != nil {
		return nil, err
	}

	ds, err := uc.delegationSeqDbRepo.GetCurrentByDelegatorUID(key)
	if err != nil {
		return nil, err
	}

	dds, err := uc.debondingDelegationSeqDbRepo.GetRecentByDelegatorUID(key, 5)
	if err != nil {
		return nil, err
	}

	return accountaggmapper.ToDetailsView(aa, ds, dds), nil
}
