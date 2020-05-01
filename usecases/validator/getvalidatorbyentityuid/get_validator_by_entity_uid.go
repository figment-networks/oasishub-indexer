package getvalidatorbyentityuid

import (
	"github.com/figment-networks/oasishub-indexer/mappers/validatoraggmapper"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatoraggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(types.PublicKey) (*validatoraggmapper.DetailsView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo               syncablerepo.DbRepo
	syncableProxyRepo            syncablerepo.ProxyRepo
	validatorAggDbRepo           validatoraggrepo.DbRepo
	validatorSeqDbRepo           validatorseqrepo.DbRepo
	delegationSeqDbRepo          delegationseqrepo.DbRepo
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	validatorAggDbRepo validatoraggrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:               syncableDbRepo,
		syncableProxyRepo:            syncableProxyRepo,
		validatorAggDbRepo:           validatorAggDbRepo,
		validatorSeqDbRepo:           validatorSeqDbRepo,
		delegationSeqDbRepo:          delegationSeqDbRepo,
		debondingDelegationSeqDbRepo: debondingDelegationSeqDbRepo,
	}
}

func (uc *useCase) Execute(key types.PublicKey) (*validatoraggmapper.DetailsView, errors.ApplicationError) {
	ea, err := uc.validatorAggDbRepo.GetByEntityUID(key)
	if err != nil {
		return nil, err
	}

	ds, err := uc.delegationSeqDbRepo.GetLastByValidatorUID(key)
	if err != nil {
		return nil, err
	}

	dds, err := uc.debondingDelegationSeqDbRepo.GetRecentByValidatorUID(key, 5)
	if err != nil {
		return nil, err
	}

	return validatoraggmapper.ToDetailsView(*ea, ds, dds), nil
}
