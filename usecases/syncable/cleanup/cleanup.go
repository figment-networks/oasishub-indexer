package cleanup

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(context.Context, int64) errors.ApplicationError
}

type useCase struct {
	syncableDbRepo   syncablerepo.DbRepo
}

func NewUseCase(syncableDbRepo syncablerepo.DbRepo) UseCase {
	return &useCase{
		syncableDbRepo:   syncableDbRepo,
	}
}

func (uc *useCase) Execute(ctx context.Context, threshold int64) errors.ApplicationError {
	h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
	if err != nil {
		return err
	}

	maxH := types.Height(h.Int64() - threshold)

	if maxH < config.FirstBlockHeight() {
		return errors.NewErrorFromMessage("nothing to cleanup", errors.CleanupError)
	}

	if err = uc.syncableDbRepo.DeletePrevByHeight(maxH); err != nil {
		return err
	}
	return nil
}
