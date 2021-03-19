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

func (a *DailyApr) calculateAPR() error {
	principalBalance := a.EscrowActiveBalance.Clone()
	if err := principalBalance.Sub(a.TotalRewards); err != nil {
		return err
	}
	if principalBalance.IsZero() {
		return errors.New("balance is zero")
	}

	a.APR = float64(a.TotalRewards.Uint64() / ((365 / 30) * principalBalance.Uint64()))

	return nil
}

func NewDailyApr(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) (DailyApr, error) {
	res := DailyApr{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
	err := res.calculateAPR()
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
