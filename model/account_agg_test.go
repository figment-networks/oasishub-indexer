package model

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"math/big"
	"testing"
	"time"
)

func Test_AccountAgg(t *testing.T) {
	model := &shared.Model{}
	agg := &shared.Aggregate{
		StartedAtHeight: int64(1),
		StartedAt:       *types.NewTimeFromTime(time.Now()),
	}

	t.Run("validation failed", func(t *testing.T) {
		m := AccountAgg{
			Model: model,
			Aggregate: agg,
		}

		if m.Valid() {
			t.Errorf("model should not be valid %+v", m)
		}
	})

	t.Run("validation success", func(t *testing.T) {
		m := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key",
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("models not equal", func(t *testing.T) {
		m1 := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key",
		}
		m2 := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key-2",
		}

		if m1.Equal(m2) {
			t.Errorf("models should not be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("models equal", func(t *testing.T) {
		m1 := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key",
		}
		m2 := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key",
		}

		if !m1.Equal(m2) {
			t.Errorf("models should be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("UpdateAggAttrs()", func(t *testing.T) {
		m1 := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key",
		}
		m2 := AccountAgg{
			Model: model,
			Aggregate: agg,

			PublicKey: "test-key",
			CurrentGeneralBalance: types.NewQuantity(big.NewInt(10)),
			CurrentGeneralNonce: uint64(1),
			CurrentEscrowActiveBalance: types.NewQuantity(big.NewInt(100)),
			CurrentEscrowActiveTotalShares: types.NewQuantity(big.NewInt(200)),
			CurrentEscrowDebondingBalance: types.NewQuantity(big.NewInt(300)),
			CurrentEscrowDebondingTotalShares: types.NewQuantity(big.NewInt(400)),
		}

		m1.UpdateAggAttrs(&m2)

		if m1.CurrentGeneralBalance.Int64() != int64(10) ||
			!m1.CurrentGeneralNonce.Equal(uint64(1)) ||
			m1.CurrentEscrowActiveBalance.Int64() != int64(100) ||
			m1.CurrentEscrowActiveTotalShares.Int64() != int64(200) ||
			m1.CurrentEscrowDebondingBalance.Int64() != int64(300) ||
			m1.CurrentEscrowDebondingTotalShares.Int64() != int64(400){
			t.Errorf("not updated %+v", m1)
		}
	})
}
