package getvalidatoruptimeforinterval

import (
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(types.PublicKey, string, string) ([]validatorseqrepo.FloatRow, errors.ApplicationError)
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

func (uc *useCase) Execute(key types.PublicKey, interval string, period string) ([]validatorseqrepo.FloatRow, errors.ApplicationError) {
	return uc.validatorSeqDbRepo.GetValidatorUptimeForInterval(key, interval, period)
}
