package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"testing"
	"time"
)

type syncableSaverStub struct{}

func (s *syncableSaverStub) Save(*syncable.Model) errors.ApplicationError {
	return nil
}

func (s *syncableSaverStub) Create(*syncable.Model) errors.ApplicationError {
	return nil
}

func Test_Sink(t *testing.T) {
	t.Run("Consume() works as expected", func(t *testing.T) {
		startH := types.Height(1)
		endH := types.Height(10)
		payload := &payload{
			StartHeight:   startH,
			EndHeight:     endH,
			CurrentHeight: startH,
			RetrievedAt:   *types.NewTimeFromTime(time.Now()),
			BlockSyncable: &syncable.Model{
				Type: syncable.BlockType,
			},
			StateSyncable: &syncable.Model{
				Type: syncable.StateType,
			},
			ValidatorsSyncable: &syncable.Model{
				Type: syncable.ValidatorsType,
			},
			TransactionsSyncable: &syncable.Model{
				Type: syncable.TransactionsType,
			},
		}
		sink := NewSink(&syncableSaverStub{}, report.Model{
			Model: &shared.Model{
				ID: 10,
			},

			StartHeight: startH,
			EndHeight:   endH,
		})
		ctx := context.Background()

		err := sink.Consume(ctx, payload)
		if err != nil {
			t.Errorf("should not return error. Err: %+v", err)
			return
		}

		if payload.BlockSyncable.ProcessedAt == nil ||
			payload.StateSyncable.ProcessedAt == nil ||
			payload.ValidatorsSyncable.ProcessedAt == nil ||
			payload.TransactionsSyncable.ProcessedAt == nil{
			t.Errorf("processedAt should be set")
		}

		if payload.BlockSyncable.ReportID == nil ||
			payload.StateSyncable.ReportID == nil ||
			payload.ValidatorsSyncable.ReportID == nil ||
			payload.TransactionsSyncable.ReportID == nil{
			t.Errorf("reportID should be set")
		}

		if sink.Count() != 1 {
			t.Errorf("should increase count. exp: %d, got: %d", 1, sink.Count())
		}
	})

}
