package apr

import (
	"errors"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DailyApr struct {
	TimeBucket          types.Time
	StartHeight         int64
	EscrowActiveBalance types.Quantity
	TotalRewards        types.Quantity
	APR                 float64
}

func calculateAPR(escrowActiveBalance, totalRewards types.Quantity) (float64, error) {
	principalBalance := escrowActiveBalance.Clone()
	if err := principalBalance.Sub(totalRewards); err != nil {
		return 0, err
	}
	if principalBalance.IsZero() {
		return 0, errors.New("balance is zero")
	}

	duration := float64(365) / float64(30)
	numerator := float64(totalRewards.Uint64()) * duration * 100
	res := numerator / float64(principalBalance.Uint64())
	return res, nil
}

func NewDailyApr(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) (DailyApr, error) {
	res := DailyApr{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
	apr, err := calculateAPR(res.EscrowActiveBalance.Clone(), res.TotalRewards.Clone())
	res.APR = apr
	return res, err
}

type MonthlyAprView struct {
	MonthInfo string  `json:"month_info"`
	AvgApr    float64 `json:"avg_apr"`
}

type MonthlyAprViewResult struct {
	Result []MonthlyAprView `json:"result"`
}

type MonthlyAprTotal struct {
	AprSum   float64
	AprCount int
}
