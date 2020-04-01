package validatordomain

type Precommit struct {
	Validated bool
	Type      int64
	Index     int64
}

func (p Precommit) Valid() bool {
	return p.Type >= 0 &&
		p.Index >= 0
}

func (p Precommit) Equal(o Precommit) bool {
	return p.Type == o.Type &&
		p.Index == o.Index
}

