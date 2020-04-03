package validatordomain

type Precommit struct {
	Validated bool  `json:"validated"`
	Type      int64 `json:"type"`
	Index     int64 `json:"index"`
}

func (p Precommit) Valid() bool {
	return p.Type >= 0 &&
		p.Index >= 0
}

func (p Precommit) Equal(o Precommit) bool {
	return p.Type == o.Type &&
		p.Index == o.Index
}
