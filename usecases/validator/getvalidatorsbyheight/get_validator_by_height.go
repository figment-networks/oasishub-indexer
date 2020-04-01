package getvalidatorsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	validatorDbRepo validatorseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	validatorDbRepo validatorseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		validatorDbRepo: validatorDbRepo,
	}
}

func (uc *useCase) Execute(height types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError) {
	txs, err := uc.validatorDbRepo.GetByHeight(height)
	if err != nil {
		return nil, err
	}

	return txs, nil
}

