package getvalidatorsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError)
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

func (uc *useCase) Execute(height *types.Height) ([]*validatordomain.ValidatorSeq, errors.ApplicationError) {
	if height == nil {
		h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
		if err != nil {
			return nil, err
		}
		height = h
	}

	es, err := uc.validatorSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	return es, nil
}

