package apr

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"math/big"
)

type RewardRate struct {
	TimeBucket          types.Time     `json:"time_bucket"`
	StartHeight         int64          `json:"start_height"`
	EscrowActiveBalance types.Quantity `json:"escrow_active_balance"`
	TotalRewards        types.Quantity `json:"total_rewards"`
	Rate                big.Float      `json:"rate"`
}

func NewRewardRate(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) RewardRate {
	res := RewardRate{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
	r := res.TotalRewards.Clone()
	rValue := new(big.Float).SetInt(r.GetBigInt())
	b := res.EscrowActiveBalance.Clone()
	bValue := new(big.Float).SetInt(b.GetBigInt())
	res.Rate = *new(big.Float).Quo(rValue, bValue)
	return res
}

type MonthlyAprView struct {
	MonthInfo string       `json:"month_info"`
	AvgApr    float64      `json:"avg_apr"`
	Dailies   []RewardRate `json:"dailies"`
}

type MonthlyAprViewResult struct {
	Result []MonthlyAprView `json:"result"`
}

type MonthlyAprTotal struct {
	MonthlyRewardRate *big.Float
	DayCount          int64
	Dailies           []RewardRate
}
