package getblockbyheight

import (
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*Response, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo       syncablerepo.DbRepo
	syncableProxyRepo    syncablerepo.ProxyRepo
	blockSeqDbRepo       blockseqrepo.DbRepo
	validatorSeqDbRepo   validatorseqrepo.DbRepo
	transactionSeqDbRepo transactionseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	blockSeqDbRepo blockseqrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	transactionSeqDbRepo transactionseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:       syncableDbRepo,
		syncableProxyRepo:    syncableProxyRepo,
		blockSeqDbRepo:       blockSeqDbRepo,
		validatorSeqDbRepo:   validatorSeqDbRepo,
		transactionSeqDbRepo: transactionSeqDbRepo,
	}
}

func (uc *useCase) Execute(height *types.Height) (*Response, errors.ApplicationError) {
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

	vs, err := uc.validatorSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}
	bs.Validators = vs

	ts, err := uc.transactionSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}
	bs.Transactions = ts

	resp := &Response{Model: bs}

	return resp, nil
}
