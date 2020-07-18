package model

import "github.com/figment-networks/oasishub-indexer/types"

const (
	SystemEventVotingPowerChange1 SystemEventKind = "voting_power_change_1"
	SystemEventVotingPowerChange2 SystemEventKind = "voting_power_change_2"
	SystemEventVotingPowerChange3 SystemEventKind = "voting_power_change_3"
	SystemEventJoinedActiveSet    SystemEventKind = "joined_active_set"
	SystemEventLeftActiveSet      SystemEventKind = "left_active_set"
	SystemEventMissedMInRow       SystemEventKind = "missed_m_in_row"
	SystemEventMissedMofN         SystemEventKind = "missed_m_of_n"
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
}

func (o SystemEvent) Update(m SystemEvent) {
	o.Height = m.Height
	o.Time = m.Time
	o.Actor = m.Actor
	o.Kind = m.Kind
}
