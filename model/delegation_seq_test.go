package model

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"math/big"
	"testing"
	"time"
)

func Test_DelegationSeq(t *testing.T) {
	chainId := "chain123"
	model := &shared.Model{}
	seq := &shared.Sequence{
		ChainId: chainId,
		Height:  int64(10),
		Time:    *types.NewTimeFromTime(time.Now()),
	}

	t.Run("validation failed", func(t *testing.T) {
		m := Report{
			Model: model,
			Sequence: seq,
		}

		if m.Valid() {
			t.Errorf("model should not be valid %+v", m)
		}
	})

	t.Run("validation success", func(t *testing.T) {
		m := Report{
			Model: model,
			Sequence: seq,

			ValidatorUID: "val-UID",
			DelegatorUID: "del-UID",
			Shares: types.NewQuantity(big.NewInt(100)),
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("models not equal", func(t *testing.T) {
		m1 := Report{
			Model: model,
			Sequence: seq,

			ValidatorUID: "val-UID",
			DelegatorUID: "del-UID",
			Shares: types.NewQuantity(big.NewInt(100)),
		}
		m2 := Report{
			Model: model,
			Sequence: seq,

			ValidatorUID: "val-UID-2",
			DelegatorUID: "del-UID-2",
			Shares: types.NewQuantity(big.NewInt(200)),
		}

		if m1.Equal(m2) {
			t.Errorf("models should not be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("models equal", func(t *testing.T) {
		m1 := Report{
			Model: model,
			Sequence: seq,

			ValidatorUID: "val-UID",
			DelegatorUID: "del-UID",
			Shares: types.NewQuantity(big.NewInt(100)),
		}
		m2 := Report{
			Model: model,
			Sequence: seq,

			ValidatorUID: "val-UID",
			DelegatorUID: "del-UID",
			Shares: types.NewQuantity(big.NewInt(100)),
		}

		if !m1.Equal(m2) {
			t.Errorf("models should be equal first: %+v; second %+v", m1, m2)
		}
	})
}


