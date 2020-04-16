package accountaggmapper

import (
	"encoding/json"
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
	"time"
)

func Test_AccountAggMapper(t *testing.T) {
	chainId := "chain123"
	model := &shared.Model{}
	rep := report.Model{
		Model:        &shared.Model{},
		StartHeight:  types.Height(1),
		EndHeight:    types.Height(10),
	}

	t.Run("ToAggregate() fails unmarshal data", func(t *testing.T) {
		s := syncable.Model{
			Model:       model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    *types.NewTimeFromTime(time.Now()),
			},

			Type:        syncable.StateType,
			Report:      rep,
			Data:        types.Jsonb{RawMessage: json.RawMessage(`{"test": 0}`)},
		}

		_, err := ToAggregate(&s)
		if err == nil {
			t.Error("data unmarshaling should fail")
		}
	})

	t.Run("ToAggregate() succeeds to unmarshal data", func(t *testing.T){
		state := `{
			"state": {
				"staking": {
					"ledger": {
						"+2ZdEh4p5JUthLrRrIz5kF4GHyzinoZ0lx6OAb/Yr/M=": {
						  "Escrow": {
							"Active": {},
							"Debonding": {},
							"StakeAccumulator": {},
							"CommissionSchedule": {}
						  },
						  "General": {
							"Nonce": "12"
						  }
						},
						"+4hqN1g5FkzRfpKvbdE4Gee2Tp5rNBIn1IMLjppVG98=": {
						  "Escrow": {
							"Active": {},
							"Debonding": {},
							"StakeAccumulator": {},
							"CommissionSchedule": {}
						  },
						  "General": {
							"Nonce": "3169"
						  }
						}
					}
				}	
			}	
		}`

		s := syncable.Model{
			Model:       model,
			Sequence: &shared.Sequence{
				ChainId: chainId,
				Height:  types.Height(10),
				Time:    *types.NewTimeFromTime(time.Now()),
			},

			Type:        syncable.StateType,
			Report:      rep,
			Data:        types.Jsonb{RawMessage: json.RawMessage(state)},
		}

		accountAggs, err := ToAggregate(&s)
		if err != nil {
			t.Error("data unmarshaling should succeed", err)
		}

		if len(accountAggs) != 2 {
			t.Errorf("there should be 2 accounts, got %d", len(accountAggs))
		}
	})
}
