package blockseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
	"time"
)

func Test_BlockSeq(t *testing.T) {
	chainId := "chain123"
	model := &shared.Model{}

	t.Run("validation failed", func(t *testing.T) {
		m := Model{
			Model: model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    time.Now(),
			},

			AppVersion: -1,
			BlockVersion: -1,
		}

		if m.Valid() {
			t.Errorf("model should not be valid %+v", m)
		}
	})

	t.Run("validation success", func(t *testing.T) {
		m := Model{
			Model: model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    time.Now(),
			},

			AppVersion: 10,
			BlockVersion: 100,
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("models not equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    time.Now(),
			},

			AppVersion: 10,
			BlockVersion: 10,
		}
		m2 := Model{
			Model: model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    time.Now(),
			},

			AppVersion: 100,
			BlockVersion: 100,
		}

		if m1.Equal(m2) {
			t.Errorf("models should not be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("models equal", func(t *testing.T) {
		tm := time.Now()
		m1 := Model{
			Model: model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    tm,
			},

			AppVersion: 10,
			BlockVersion: 10,
		}
		m2 := Model{
			Model: model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    tm,
			},

			AppVersion: 10,
			BlockVersion: 10,
		}

		if !m1.Equal(m2) {
			t.Errorf("models should be equal first: %+v; second %+v", m1, m2)
		}
	})
}
