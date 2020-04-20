package validatorseqmapper

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

func Test_StakingSeqMapper(t *testing.T) {
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
	validatorsFixture := fixtures.Load("validators.json")
	blockFixture := fixtures.Load("block.json")
	stateFixture := fixtures.Load("state.json")

	t.Run("ToSequence()() fails unmarshal data", func(t *testing.T) {
		vs := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.ValidatorsType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(`{"test": 0}`)},
		}

		bs := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.BlockType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(`{"test": 0}`)},
		}

		ss := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.StateType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(`{"test": 0}`)},
		}

		_, err := ToSequence(vs, bs, ss)
		if err == nil {
			t.Error("data unmarshaling should fail")
		}
	})

	t.Run("ToSequence()() succeeds to unmarshal data", func(t *testing.T) {
		vs := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.ValidatorsType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(validatorsFixture)},
		}

		bs := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.BlockType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(blockFixture)},
		}

		ss := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.StateType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(stateFixture)},
		}

		validatorSeqs, err := ToSequence(vs, bs, ss)
		if err != nil {
			t.Error("data unmarshaling should succeed", err)
		}

		if len(validatorSeqs) == 0 {
			t.Error("there should be transactions")
		}
	})
}
