package entityagg

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
	"time"
)

func Test_EntityAgg(t *testing.T) {
	model := &shared.Model{}
	agg := &shared.Aggregate{
		StartedAtHeight: types.Height(1),
		StartedAt:       *types.NewTimeFromTime(time.Now()),
	}

	t.Run("validation failed", func(t *testing.T) {
		m := Model{
			Model: model,
			Aggregate: agg,
		}

		if m.Valid() {
			t.Errorf("model should not be valid %+v", m)
		}
	})

	t.Run("validation success", func(t *testing.T) {
		m := Model{
			Model: model,
			Aggregate: agg,

			EntityUID: "test-UID",
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("models not equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Aggregate: agg,

			EntityUID: "test-UID",
		}
		m2 := Model{
			Model: model,
			Aggregate: agg,

			EntityUID: "test-UID-2",
		}

		if m1.Equal(m2) {
			t.Errorf("models should not be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("models equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Aggregate: agg,

			EntityUID: "test-UID",
		}
		m2 := Model{
			Model: model,
			Aggregate: agg,

			EntityUID: "test-UID",
		}

		if !m1.Equal(m2) {
			t.Errorf("models should be equal first: %+v; second %+v", m1, m2)
		}
	})
}

