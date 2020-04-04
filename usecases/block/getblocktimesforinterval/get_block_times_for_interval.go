package getblocktimesforinterval

import (
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(string, string) ([]blockseqrepo.Row, errors.ApplicationError)
}

type useCase struct {
	blockDbRepo blockseqrepo.DbRepo
}

func NewUseCase(blockDbRepo blockseqrepo.DbRepo) UseCase {
	return &useCase{
		blockDbRepo:   blockDbRepo,
	}
}

func (uc *useCase) Execute(interval string, period string) ([]blockseqrepo.Row, errors.ApplicationError) {
	return uc.blockDbRepo.GetAvgBlockTimesForInterval(interval, period)
}
