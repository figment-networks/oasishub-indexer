package getaccountbypublickey

import (
	"github.com/figment-networks/oasishub-indexer/domain/accountdomain"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(key types.PublicKey) (*accountdomain.AccountAgg, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	accountSeqDbRepo       accountaggrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	accountSeqDbRepo accountaggrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		accountSeqDbRepo:       accountSeqDbRepo,
	}
}

func (uc *useCase) Execute(k types.PublicKey) (*accountdomain.AccountAgg, errors.ApplicationError) {
	aa, err := uc.accountSeqDbRepo.GetByPublicKey(k)
	if err != nil {
		return nil, err
	}

	return aa, nil
}

