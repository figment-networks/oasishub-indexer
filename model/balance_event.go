package model

import "github.com/figment-networks/oasishub-indexer/types"

const (
	Commission     BalanceEventKind = "commission"
	SlashActive    BalanceEventKind = "slash_active"
	SlashDebonding BalanceEventKind = "slash_debonding"
	Reward         BalanceEventKind = "reward"
)

type BalanceEventKind string

func (b BalanceEventKind) String() string {
	return string(b)
}

type BalanceEvent struct {
	*Model

	Height        int64            `json:"height"`
	Address       string           `json:"address"`
	EscrowAddress string           `json:"escrow_addr"`
	Amount        types.Quantity   `json:"amount"`
	Kind          BalanceEventKind `json:"kind"`
}

func (b BalanceEvent) Update(m BalanceEvent) {
	b.Height = m.Height
	b.Address = m.Address
	b.EscrowAddress = m.EscrowAddress
	b.Amount = m.Amount
	b.Kind = m.Kind
}
