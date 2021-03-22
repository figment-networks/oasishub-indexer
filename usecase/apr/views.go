package apr

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"math/big"
)

type DailyApr struct {
	TimeBucket          types.Time
	StartHeight         int64
	EscrowActiveBalance types.Quantity
	TotalRewards        types.Quantity
	rate                big.Float
}

func NewDailyApr(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) DailyApr {
	res := DailyApr{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
	res.rate = *new(big.Float)
	res.rate.SetFloat64(float64(res.TotalRewards.Int64()) / float64(res.EscrowActiveBalance.Int64()))
	return res
}

type MonthlyAprView struct {
	MonthInfo string  `json:"month_info"`
	AvgApr    float64 `json:"avg_apr"`
}

type MonthlyAprViewResult struct {
	Result []MonthlyAprView `json:"result"`
}

type MonthlyAprTotal struct {
	MonthlyRewardRate *big.Float
	DayCount          int64
}
