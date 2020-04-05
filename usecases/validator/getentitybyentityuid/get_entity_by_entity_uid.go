package getentitybyentityuid

import (
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(types.PublicKey) (*Response, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo               syncablerepo.DbRepo
	syncableProxyRepo            syncablerepo.ProxyRepo
	entityAggDbRepo              entityaggrepo.DbRepo
	validatorSeqDbRepo           validatorseqrepo.DbRepo
	delegationSeqDbRepo          delegationseqrepo.DbRepo
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	entityAggDbRepo entityaggrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:               syncableDbRepo,
		syncableProxyRepo:            syncableProxyRepo,
		entityAggDbRepo:              entityAggDbRepo,
		validatorSeqDbRepo:           validatorSeqDbRepo,
		delegationSeqDbRepo:          delegationSeqDbRepo,
		debondingDelegationSeqDbRepo: debondingDelegationSeqDbRepo,
	}
}

func (uc *useCase) Execute(key types.PublicKey) (*Response, errors.ApplicationError) {
	ea, err := uc.entityAggDbRepo.GetByEntityUID(key)
	if err != nil {
		return nil, err
	}

	resp := &Response{Model: ea}

	h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
	if err != nil {
		return nil, err
	}
	resp.LastHeight = *h

	tv, err := uc.validatorSeqDbRepo.GetTotalValidatedByEntityUID(key)
	if err != nil {
		return nil, err
	}
	resp.TotalValidated = *tv

	tm, err := uc.validatorSeqDbRepo.GetTotalMissedByEntityUID(key)
	if err != nil {
		return nil, err
	}
	resp.TotalMissed = *tm

	tp, err := uc.validatorSeqDbRepo.GetTotalProposedByEntityUID(key)
	if err != nil {
		return nil, err
	}
	resp.TotalProposed = *tp

	ds, err := uc.delegationSeqDbRepo.GetLastByValidatorUID(key)
	if err != nil {
		return nil, err
	}
	resp.LastDelegations = ds

	dds, err := uc.debondingDelegationSeqDbRepo.GetRecentByValidatorUID(key, 5)
	if err != nil {
		return nil, err
	}
	resp.RecentDebondingDelegations = dds

	return resp, nil
}
