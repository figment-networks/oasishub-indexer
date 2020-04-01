package gettransactionsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/transactiondomain"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height types.Height) ([]*transactiondomain.TransactionSeq, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	transactionDbRepo transactionseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	transactionDbRepo transactionseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		transactionDbRepo: transactionDbRepo,
	}
}

func (uc *useCase) Execute(height types.Height) ([]*transactiondomain.TransactionSeq, errors.ApplicationError) {
	txs, err := uc.transactionDbRepo.GetByHeight(height)
	if err != nil {
		return nil, err
	}

	return txs, nil
}
