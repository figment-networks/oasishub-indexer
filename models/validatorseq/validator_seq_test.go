package validatorseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
	"time"
)

func Test_TransactionSeq(t *testing.T) {
	chainId := "chain123"
	model := &shared.Model{}
	seq := &shared.Sequence{
		ChainId: chainId,
		Height:  types.Height(10),
		Time:    *types.NewTimeFromTime(time.Now()),
	}

	t.Run("validation failed", func(t *testing.T) {
		m := Model{
			Model: model,
			Sequence: seq,
		}

		if m.Valid() {
			t.Errorf("model should not be valid %+v", m)
		}
	})

	t.Run("validation success", func(t *testing.T) {
		m := Model{
			Model: model,
			Sequence: seq,

			EntityUID: types.PublicKey("entity-uid-test"),
			NodeUID: types.PublicKey("node-uid-test"),
			VotingPower: 10,
			TotalShares: types.NewQuantityFromInt64(int64(100)),
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("models not equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Sequence: seq,

			EntityUID: types.PublicKey("entity-uid-test"),
			NodeUID: types.PublicKey("node-uid-test"),
			VotingPower: 10,
			TotalShares: types.NewQuantityFromInt64(int64(100)),
		}
		m2 := Model{
			Model: model,
			Sequence: seq,

			EntityUID: types.PublicKey("entity-uid-test-2"),
			NodeUID: types.PublicKey("node-uid-test-2"),
			VotingPower: 20,
			TotalShares: types.NewQuantityFromInt64(int64(200)),
		}

		if m1.Equal(m2) {
			t.Errorf("models should not be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("models equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Sequence: seq,

			EntityUID: types.PublicKey("entity-uid-test"),
			NodeUID: types.PublicKey("node-uid-test"),
			VotingPower: 10,
			TotalShares: types.NewQuantityFromInt64(int64(100)),
		}
		m2 := Model{
			Model: model,
			Sequence: seq,

			EntityUID: types.PublicKey("entity-uid-test"),
			NodeUID: types.PublicKey("node-uid-test"),
			VotingPower: 10,
			TotalShares: types.NewQuantityFromInt64(int64(100)),
		}

		if !m1.Equal(m2) {
			t.Errorf("models should be equal first: %+v; second %+v", m1, m2)
		}
	})
}


