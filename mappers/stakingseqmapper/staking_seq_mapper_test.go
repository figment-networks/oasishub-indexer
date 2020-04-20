package stakingseqmapper

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
	stateFixture := fixtures.Load("state.json")

	t.Run("ToSequence()() fails unmarshal data", func(t *testing.T) {
		s := syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type:   syncable.StateType,
			Report: rep,
			Data:   types.Jsonb{RawMessage: json.RawMessage(`{"test": 0}`)},
		}

		_, err := ToSequence(s)
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

		stakingSeq, err := ToSequence(s)
		if err != nil {
			t.Error("data unmarshaling should succeed", err)
		}

		exp := types.NewQuantityFromBytes([]byte("iscjBInoAAA="))
		if stakingSeq.TotalSupply.Equals(exp) {
			t.Errorf("wrong total supply, exp: %v, got: %v", exp, stakingSeq.TotalSupply)
		}

		exp2 := types.NewQuantityFromBytes([]byte("bwIHGjdarc8="))
		if stakingSeq.CommonPool.Equals(exp2) {
			t.Errorf("wrong common pool, exp: %v, got: %v", exp2, stakingSeq.CommonPool)
		}

		exp3 := uint64(10)
		if stakingSeq.DebondingInterval != exp3 {
			t.Errorf("wrong debonding interval, exp: %d, got: %d", exp3, stakingSeq.DebondingInterval)
		}

		exp4 := types.NewQuantityFromBytes([]byte("AlQL5AA="))
		if stakingSeq.CommonPool.Equals(exp4) {
			t.Errorf("wrong common pool, exp: %v, got: %v", exp4, stakingSeq.MinDelegationAmount)
		}
	})
}

