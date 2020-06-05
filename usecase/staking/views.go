package staking

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DetailsView struct {
	TotalSupply         types.Quantity `json:"total_supply"`
	CommonPool          types.Quantity `json:"common_pool"`
	DebondingInterval   uint64         `json:"debonding_interval"`
	MinDelegationAmount types.Quantity `json:"min_delegation_amount"`
}

func ToDetailsView(rawStaking *statepb.Staking) *DetailsView {
	return &DetailsView{
		TotalSupply:         types.NewQuantityFromBytes(rawStaking.GetTotalSupply()),
		CommonPool:          types.NewQuantityFromBytes(rawStaking.GetCommonPool()),
		DebondingInterval:   rawStaking.GetParameters().GetDebondingInterval(),
		MinDelegationAmount: types.NewQuantityFromBytes(rawStaking.GetParameters().GetMinDelegationAmount()),
	}
}
