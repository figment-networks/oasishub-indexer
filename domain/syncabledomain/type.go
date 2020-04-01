package syncabledomain

const (
	BlockType        = "block"
	StateType        = "state"
	ValidatorsType   = "validators"
	TransactionsType = "transactions"
)

var Types = []Type{BlockType, StateType, ValidatorsType, TransactionsType}

type Type string

func (t Type) Valid() bool {
	return t == BlockType ||
		t == StateType ||
		t == ValidatorsType ||
		t == TransactionsType
}
