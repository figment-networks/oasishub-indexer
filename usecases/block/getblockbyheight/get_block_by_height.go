package getblockbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/blockdomain"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*blockdomain.BlockSeq, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	blockSeqDbRepo       blockseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	blockSeqDbRepo blockseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		blockSeqDbRepo:       blockSeqDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*blockdomain.BlockSeq, errors.ApplicationError) {
	if height == nil {
		h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
		if err != nil {
			return nil, err
		}
		height = h
	}

	bs, err := uc.blockSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
