package apr

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"math/big"
)

type DailyApr struct {
	TimeBucket          types.Time     `json:"time_bucket"`
	StartHeight         int64          `json:"start_height"`
	EscrowActiveBalance types.Quantity `json:"escrow_active_balance"`
	TotalRewards        types.Quantity `json:"total_rewards"`
	Rate                big.Float      `json:"rate"`
}

func NewDailyApr(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) DailyApr {
	res := DailyApr{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
	res.Rate = *new(big.Float)
	res.Rate.SetFloat64(float64(res.TotalRewards.Int64()) / float64(res.EscrowActiveBalance.Int64()))
	return res
}

type MonthlyAprView struct {
	MonthInfo string     `json:"month_info"`
	AvgApr    float64    `json:"avg_apr"`
	Dailies   []DailyApr `json:"dailies"`
}

type MonthlyAprViewResult struct {
	Result []MonthlyAprView `json:"result"`
}

type MonthlyAprTotal struct {
	MonthlyRewardRate *big.Float
	DayCount          int64
	Dailies           []DailyApr
}
