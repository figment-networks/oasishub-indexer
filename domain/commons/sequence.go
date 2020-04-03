package commons

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"time"
)

type Sequence struct {
	ChainId string       `json:"chain_id"`
	Height  types.Height `json:"height"`
	Time    time.Time    `json:"time"`
}

type SequenceProps struct {
	ChainId string
	Height  types.Height
	Time    time.Time
}

func NewSequence(props SequenceProps) *Sequence {
	return &Sequence{
		ChainId: props.ChainId,
		Height:  props.Height,
		Time:    props.Time,
	}
}

func (s *Sequence) Valid() bool {
	return s.ChainId != "" &&
		s.Height.Valid() &&
		!s.Time.IsZero()
}

func (s *Sequence) Equal(m Sequence) bool {
	return s.ChainId == m.ChainId &&
		s.Height.Equal(m.Height) &&
		s.Time.Equal(m.Time)
}
