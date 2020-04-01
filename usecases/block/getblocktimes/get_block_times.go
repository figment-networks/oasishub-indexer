package getblocktimes

import (
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
)

type UseCase interface {
	Execute(int64) blockseqrepo.Result
}

type useCase struct {
	blockDbRepo blockseqrepo.DbRepo
}

func NewUseCase(blockDbRepo blockseqrepo.DbRepo) UseCase {
	return &useCase{
		blockDbRepo:   blockDbRepo,
	}
}

func (uc *useCase) Execute(limit int64) blockseqrepo.Result {
	return uc.blockDbRepo.GetAvgBlockTimesForRecentBlocks(limit)
}
