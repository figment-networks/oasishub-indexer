package gettotalsharesforinterval

import (
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(string, string) ([]validatorseqrepo.Row, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	validatorSeqDbRepo validatorseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		validatorSeqDbRepo: validatorSeqDbRepo,
	}
}

func (uc *useCase) Execute(interval string, period string) ([]validatorseqrepo.Row, errors.ApplicationError) {
	return uc.validatorSeqDbRepo.GetTotalSharesForInterval(interval, period)
}
