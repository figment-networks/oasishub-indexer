package staking

import (
	"github.com/figment-networks/oasishub-indexer/model"
)

type DetailsView struct {
	*model.StakingSeq
}

func ToDetailsView(s *model.StakingSeq) *DetailsView {
	return &DetailsView{
		StakingSeq: s,
	}
}
