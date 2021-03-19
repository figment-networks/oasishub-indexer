package apr

import (
	"errors"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

type MonthlyAprView struct {
	TimeBucket          types.Time     `json:"time_bucket"`
	StartHeight         int64          `json:"height"`
	EscrowActiveBalance types.Quantity `json:"escrow_active_balance"`
	TotalRewards        types.Quantity `json:"total_rewards"`
	APR                 float64        `json:"apr"`
}

func (a *MonthlyAprView) calculateAPR() error {
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

func NewMonthlyAprView(summary model.BalanceSummary, rawAccount *accountpb.GetByAddressResponse) (MonthlyAprView, error) {
	res := MonthlyAprView{
		TimeBucket:          summary.TimeBucket,
		StartHeight:         summary.StartHeight,
		EscrowActiveBalance: types.NewQuantityFromBytes(rawAccount.GetAccount().GetEscrow().GetActive().GetBalance()),
		TotalRewards:        summary.TotalRewards,
	}
	err := res.calculateAPR()
	return res, err
}

type MonthlyAprViewResult struct {
	Result []MonthlyAprView `json:"result"`
}
