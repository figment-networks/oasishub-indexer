package apr

import (
	"fmt"
	"math/big"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

var (
	daysInYear   = big.NewFloat(365)
	decPrecision = 4
)

type DailyApr struct {
	TimeBucket   string         `json:"time_bucket"`
	Bonded       types.Quantity `json:"bonded"`
	TotalRewards types.Quantity `json:"total_rewards"`
	APR          string         `json:"apr"`
	Validator    string         `json:"validator"`
}

func toAPRView(summaries []model.ValidatorSummary, rewardLookup map[string]model.BalanceSummary, delegationLookup map[string]*delegationpb.GetByAddressResponse) (res []DailyApr, err error) {
	for _, s := range summaries {
		rewardSeq, ok := rewardLookup[fmt.Sprintf("%s.%s", s.Address, s.TimeBucket.Format(timeFormat))]
		if !ok {
			return res, fmt.Errorf("missing reward for address %s and time %s", s.Address, s.TimeBucket.Format(timeFormat))
		}

		delegation, ok := delegationLookup[fmt.Sprintf("%s.%d", rewardSeq.Address, rewardSeq.StartHeight)]
		if !ok {
			return res, fmt.Errorf("missing delegation for address %s and height %d", rewardSeq.Address, rewardSeq.StartHeight)
		}

		stake, err := getStakedBalance(s, delegation)
		if err != nil {
			return res, err
		}

		res = append(res, dailyAPR(rewardSeq, stake))
	}

	return res, nil

}

func dailyAPR(rewardSeq model.BalanceSummary, stake types.Quantity) DailyApr {
	r := rewardSeq.TotalRewards.Clone()
	rValue := new(big.Float).SetInt(&r.Int)
	b := stake.Clone()
	bValue := new(big.Float).SetInt(&b.Int)
	apr := new(big.Float).Quo(rValue, bValue)
	apr = apr.Mul(apr, daysInYear)

	return DailyApr{
		TimeBucket:   rewardSeq.TimeBucket.Format(timeFormat),
		Bonded:       stake,
		TotalRewards: rewardSeq.TotalRewards,
		APR:          apr.Text('f', decPrecision),
		Validator:    rewardSeq.EscrowAddress,
	}
}

func getStakedBalance(validator model.ValidatorSummary, delegations *delegationpb.GetByAddressResponse) (types.Quantity, error) {
	d, ok := delegations.GetDelegations()[validator.Address]
	if !ok {
		return types.Quantity{}, fmt.Errorf("missing delegation for validator %s", validator.Address)
	}

	// value_per_share = total_base_units / total_shares
	// delegated_balance = value_per_share * total_delegated_shares
	// rewrite to multiply first: delegated_balance = total_base_units * total_delegated_shares / total_shares
	balance := types.NewQuantityFromBytes(d.GetShares())

	err := balance.Mul(validator.ActiveEscrowBalanceAvg)
	if !ok {
		return types.Quantity{}, err
	}

	err = balance.Quo(validator.TotalSharesAvg)
	return balance, err
}
