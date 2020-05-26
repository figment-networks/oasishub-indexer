package model

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
)

func Test_Report(t *testing.T) {
	model := &shared.Model{}

	t.Run("validation success", func(t *testing.T) {
		m := Model{
			Model: model,

			StartHeight: int64(10),
			EndHeight: int64(20),
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("Complete()", func(t *testing.T) {
		m := Model{
			Model: model,

			StartHeight: int64(10),
			EndHeight: int64(20),
		}

		m.Complete(int64(10), int64(5), nil, nil)

		if *m.SuccessCount != int64(10) ||
			*m.ErrorCount != int64(5) ||
			m.ErrorMsg != nil ||
			m.Details.RawMessage != nil {
			t.Errorf("values not updated %+v", m)
		}
	})
}


