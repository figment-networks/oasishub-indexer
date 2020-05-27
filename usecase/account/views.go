package account

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
)

type DetailsView struct {
	*accountpb.Account

	CurrentDelegations         []model.DelegationSeq          `json:"current_delegations"`
	RecentDebondingDelegations []model.DebondingDelegationSeq `json:"recent_debonding_delegations"`
}

func ToDetailsView(rawAccount *accountpb.Account, accountAgg *model.AccountAgg, ds []model.DelegationSeq, dds []model.DebondingDelegationSeq) *DetailsView {
	return &DetailsView{
		Account: rawAccount,

		CurrentDelegations:         ds,
		RecentDebondingDelegations: dds,
	}
}
