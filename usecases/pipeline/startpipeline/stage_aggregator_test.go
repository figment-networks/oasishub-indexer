package startpipeline

import (
	"context"
	"encoding/json"
	"github.com/figment-networks/oasishub-indexer/fixtures"
	mock_accountaggrepo "github.com/figment-networks/oasishub-indexer/mock/repos/accountaggrepo"
	mock_validatoraggrepo "github.com/figment-networks/oasishub-indexer/mock/repos/validatoraggrepo"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func Test_StageAggregator(t *testing.T) {
	blockFixture := fixtures.Load("block.json")
	stateFixture := fixtures.Load("state.json")
	validatorsFixture := fixtures.Load("validators.json")
	transactionsFixture := fixtures.Load("transactions.json")

	startH := types.Height(1)
	endH := types.Height(10)
	chainId := "chain123"
	model := &shared.Model{}
	sequence := &shared.Sequence{
		ChainId: chainId,
		Height:  types.Height(10),
		Time:    *types.NewTimeFromTime(time.Now()),
	}
	pld := &payload{
		StartHeight:   startH,
		EndHeight:     endH,
		CurrentHeight: startH,
		RetrievedAt:   *types.NewTimeFromTime(time.Now()),
		BlockSyncable: &syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type: syncable.BlockType,
			Data: types.Jsonb{RawMessage: json.RawMessage(blockFixture)},
		},
		StateSyncable: &syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type: syncable.StateType,
			Data: types.Jsonb{RawMessage: json.RawMessage(stateFixture)},
		},
		ValidatorsSyncable: &syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type: syncable.ValidatorsType,
			Data: types.Jsonb{RawMessage: json.RawMessage(validatorsFixture)},
		},
		TransactionsSyncable: &syncable.Model{
			Model:    model,
			Sequence: sequence,

			Type: syncable.TransactionsType,
			Data: types.Jsonb{RawMessage: json.RawMessage(transactionsFixture)},
		},
	}
	ctx := context.Background()

	t.Run("Process() works as expected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		accountDbRepo := mock_accountaggrepo.NewMockDbRepo(ctrl)
		accountDbRepo.EXPECT().GetByPublicKey(gomock.Any()).Return(nil, errors.NewErrorFromMessage("not found", errors.NotFoundError)).MinTimes(1)
		accountDbRepo.EXPECT().Create(gomock.Any()).Return(nil).MinTimes(1)

		entityDbRepo := mock_validatoraggrepo.NewMockDbRepo(ctrl)
		entityDbRepo.EXPECT().GetByEntityUID(gomock.Any()).Return(nil, errors.NewErrorFromMessage("not found", errors.NotFoundError)).MinTimes(1)
		entityDbRepo.EXPECT().Create(gomock.Any()).Return(nil).MinTimes(1)

		aggregator := NewAggregator(accountDbRepo, entityDbRepo)

		returnedPayload, err := aggregator.Process(ctx, pld)
		if err != nil {
			t.Errorf("should not return error. Err: %+v", err)
			return
		}
		p := returnedPayload.(*payload)

		if p.NewAggregatedAccounts == nil {
			t.Errorf("new accounts should be aggregated")
		}
		if p.NewAggregatedValidators == nil {
			t.Errorf("new entities should be aggregated")
		}
	})
}
