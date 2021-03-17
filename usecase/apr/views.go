package apr

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

type MonthlyAprView struct {
	TimeBucket          types.Time     `json:"time_bucket"`
	StartHeight         int64          `json:"height"`
	EscrowActiveBalance types.Quantity `json:"escrow_active_balance"`
	TotalRewards        types.Quantity `json:"total_rewards"`
}

func NewMonthlyAprView(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) *MonthlyAprView {
	return &MonthlyAprView{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
}

type MonthlyAprViewResult struct {
	Result []MonthlyAprView `json:"result"`
}
