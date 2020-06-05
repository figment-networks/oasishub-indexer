package account

import (
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
