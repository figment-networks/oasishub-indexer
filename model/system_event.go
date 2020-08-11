package model

import "github.com/figment-networks/oasishub-indexer/types"

const (
	SystemEventActiveEscrowBalanceChange1 SystemEventKind = "active_escrow_balance_change_1"
	SystemEventActiveEscrowBalanceChange2 SystemEventKind = "active_escrow_balance_change_2"
	SystemEventActiveEscrowBalanceChange3 SystemEventKind = "active_escrow_balance_change_3"
	SystemEventCommissionChange1          SystemEventKind = "commission_change_1"
	SystemEventCommissionChange2          SystemEventKind = "commission_change_2"
	SystemEventCommissionChange3          SystemEventKind = "commission_change_3"
	SystemEventJoinedActiveSet            SystemEventKind = "joined_active_set"
	SystemEventLeftActiveSet              SystemEventKind = "left_active_set"
	SystemEventMissedMInRow               SystemEventKind = "missed_m_in_row"
	SystemEventMissedMofN                 SystemEventKind = "missed_m_of_n"
)

type SystemEventKind string

func (o SystemEventKind) String() string {
	return string(o)
}

type SystemEvent struct {
	*Model

	Height int64           `json:"height"`
	Time   types.Time      `json:"time"`
	Actor  string          `json:"actor"`
	Kind   SystemEventKind `json:"kind"`
	Data   types.Jsonb     `json:"data"`
}

func (o SystemEvent) Update(m SystemEvent) {
	o.Height = m.Height
	o.Time = m.Time
	o.Actor = m.Actor
	o.Kind = m.Kind
	o.Data = m.Data
}
