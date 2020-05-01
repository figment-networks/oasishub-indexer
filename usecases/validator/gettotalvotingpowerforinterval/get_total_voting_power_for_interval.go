package gettotalvotingpowerforinterval

import (
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type UseCase interface {
	Execute(string, string) ([]validatorseqrepo.FloatRow, errors.ApplicationError)
}

type useCase struct {
	validatorSeqDbRepo validatorseqrepo.DbRepo
}

func NewUseCase(
	validatorSeqDbRepo validatorseqrepo.DbRepo,
) UseCase {
	return &useCase{
		validatorSeqDbRepo: validatorSeqDbRepo,
	}
}

func (uc *useCase) Execute(interval string, period string) ([]validatorseqrepo.FloatRow, errors.ApplicationError) {
	return uc.validatorSeqDbRepo.GetTotalVotingPowerForInterval(interval, period)
}
