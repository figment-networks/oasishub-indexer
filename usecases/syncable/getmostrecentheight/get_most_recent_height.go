package getmostrecentheight

import (
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute() (*types.Height, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo syncablerepo.DbRepo
}

func NewUseCase(syncableDbRepo syncablerepo.DbRepo) UseCase {
	return &useCase{
		syncableDbRepo: syncableDbRepo,
	}
}

func (uc *useCase) Execute() (*types.Height, errors.ApplicationError) {
	h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
	if err != nil {
		return nil, err
	}
	return h, nil
}
