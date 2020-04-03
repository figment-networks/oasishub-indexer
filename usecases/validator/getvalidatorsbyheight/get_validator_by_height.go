package getvalidatorsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(height *types.Height) (*Response, errors.ApplicationError)
}

type useCase struct {
	syncableDbRepo    syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo
	validatorSeqDbRepo validatorseqrepo.DbRepo
	delegationSeqDbRepo delegationseqrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:    syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,
		validatorSeqDbRepo: validatorSeqDbRepo,
		delegationSeqDbRepo: delegationSeqDbRepo,
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

	vs, err := uc.validatorSeqDbRepo.GetByHeight(*height)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Validators: vs,
	}
	return resp, nil
}

