package getvalidatorsforminheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/validatoraggmapper"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatoraggrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

var _ UseCase = (*useCase)(nil)

type UseCase interface {
	Execute(height *types.Height) (*validatoraggmapper.ListView, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	validatorAggDbRepo validatoraggrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	validatorAggDbRepo validatoraggrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		validatorAggDbRepo: validatorAggDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*validatoraggmapper.ListView, errors.ApplicationError) {
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

	vas, err := uc.validatorAggDbRepo.GetAllForHeightGreaterThan(*height)
	if err != nil {
		return nil, err
	}

	return validatoraggmapper.ToListView(vas), nil
}

