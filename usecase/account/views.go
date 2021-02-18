package account

import (
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DetailsView struct {
	GeneralBalance             types.Quantity `json:"general_balance"`
	GeneralNonce               uint64         `json:"general_nonce"`
	EscrowActiveBalance        types.Quantity `json:"escrow_active_balance"`
	EscrowActiveTotalShares    types.Quantity `json:"escrow_active_total_shares"`
	EscrowDebondingBalance     types.Quantity `json:"escrow_debonding_balance"`
	EscrowDebondingTotalShares types.Quantity `json:"escrow_debonding_total_shares"`
}

func ToDetailsView(rawAccount *accountpb.Account) *DetailsView {
	return &DetailsView{
		GeneralBalance:             types.NewQuantityFromBytes(rawAccount.GetGeneral().GetBalance()),
		GeneralNonce:               rawAccount.GetGeneral().GetNonce(),
		EscrowActiveBalance:        types.NewQuantityFromBytes(rawAccount.GetEscrow().GetActive().GetBalance()),
		EscrowActiveTotalShares:    types.NewQuantityFromBytes(rawAccount.GetEscrow().GetActive().GetBalance()),
		EscrowDebondingBalance:     types.NewQuantityFromBytes(rawAccount.GetEscrow().GetDebonding().GetBalance()),
		EscrowDebondingTotalShares: types.NewQuantityFromBytes(rawAccount.GetEscrow().GetDebonding().GetTotalShares()),
	}
}

type DailyBalanceViewResult struct {
	Result []DailyBalanceView `json:"result"`
}
type DailyBalanceView struct {
	TimeBucket time.Time   `json:"time_bucket"`
	Balance    DetailsView `json:"balance"`
}

func toDailyBalanceViewResult(data []dataRow) DailyBalanceViewResult {
	var summaries []DailyBalanceView
	for _, row := range data {
		summaries = append(summaries, DailyBalanceView{
			TimeBucket: row.dayStart,
			Balance: DetailsView{
				GeneralBalance:             types.NewQuantityFromBytes(row.account.GetGeneral().GetBalance()),
				GeneralNonce:               row.account.GetGeneral().GetNonce(),
				EscrowActiveBalance:        types.NewQuantityFromBytes(row.account.GetEscrow().GetActive().GetBalance()),
				EscrowActiveTotalShares:    types.NewQuantityFromBytes(row.account.GetEscrow().GetActive().GetBalance()),
				EscrowDebondingBalance:     types.NewQuantityFromBytes(row.account.GetEscrow().GetDebonding().GetBalance()),
				EscrowDebondingTotalShares: types.NewQuantityFromBytes(row.account.GetEscrow().GetDebonding().GetTotalShares()),
			},
		})
	}
	return DailyBalanceViewResult{Result: summaries}
}
