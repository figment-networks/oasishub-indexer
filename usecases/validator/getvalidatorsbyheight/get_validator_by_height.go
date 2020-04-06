package getvalidatorsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/validatorseqmapper"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

var _ UseCase = (*useCase)(nil)

type UseCase interface {
	Execute(height *types.Height) (*validatorseqmapper.ListView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	validatorSeqDbRepo validatorseqrepo.DbRepo
	delegationSeqDbRepo delegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		validatorSeqDbRepo: validatorSeqDbRepo,
		delegationSeqDbRepo: delegationSeqDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*validatorseqmapper.ListView, errors.ApplicationError) {
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

	vs, err := uc.validatorSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return validatorseqmapper.ToListView(vs), nil
}

