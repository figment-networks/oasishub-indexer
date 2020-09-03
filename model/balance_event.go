package model

import "github.com/figment-networks/oasishub-indexer/types"

const (
	Commission BalanceEventKind = "commission"
	Slash      BalanceEventKind = "slash"
	Reward     BalanceEventKind = "reward"
)

type BalanceEventKind string

func (o BalanceEventKind) String() string {
	return string(o)
}

type BalanceEvent struct {
	*Model

	Height        int64            `json:"height"`
	Address       string           `json:"address"`
	EscrowAddress string           `json:"escrow_addr"`
	Amount        types.Quantity   `json:"amount"`
	Kind          BalanceEventKind `json:"kind"`
}

// func (o BalanceEvent) Update(m BalanceEvent) {
// 	o.Height = m.Height
// 	o.Time = m.Time
// 	o.Actor = m.Actor
// 	o.Kind = m.Kind
// 	o.Data = m.Data
// }

// stakeForShares
// single = 1*600/333 = 1.8 = balance/total_shares
// share value? = 1.8 base_units/share
// prev share value =
// 1. reverse com deposit = 600-60=540, 60*.55=33 333-33=300 = shares
// 2. reverse reward, 540-240 =300 = balance
// calc share value 300/300 = 1 balance/total_shares
// increased share value = 1.8 - 1 = 0.8

// rewards for del1 = 100*.8 = 80
// rewards for del2 = 200*.8 = 160
// total = 240!!
