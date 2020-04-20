package delegationseqmapper

import (
	"encoding/json"
	"github.com/figment-networks/oasishub-indexer/fixtures"
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
	"time"
)

func Test_DelegationSeqMapper(t *testing.T) {
	chainId := "chain123"
	model := &shared.Model{}
	sequence := &shared.Sequence{
		ChainId: chainId,
		Height:  types.Height(10),
		Time:    *types.NewTimeFromTime(time.Now()),
	}
	rep := report.Model{
		Model:       &shared.Model{},
		StartHeight: types.Height(1),
		EndHeight:   types.Height(10),
	}
	stateFixture := fixtures.Load("state.json")

	t.Run("ToSequence()() fails unmarshal data", func(t *testing.T) {
		s := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.StateType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(`{"test": 0}`)},
		}

		_, err := ToSequence(&s)
		if err == nil {
			t.Error("data unmarshaling should fail")
		}
	})

	t.Run("ToSequence()() succeeds to unmarshal data", func(t *testing.T) {
		s := syncable.Model{
			Model: model,
			Sequence: sequence,

			Type:   syncable.StateType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(stateFixture)},
		}

		delegationSeqs, err := ToSequence(&s)
		if err != nil {
			t.Error("data unmarshaling should succeed", err)
		}

		if len(delegationSeqs) == 0 {
			t.Error("there should be accounts")
		}
	})
}
