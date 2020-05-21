package startpipeline

import (
	"context"
	"encoding/json"
	"github.com/figment-networks/oasishub-indexer/fixtures"
	mock_blockseqrepo "github.com/figment-networks/oasishub-indexer/mock/repos/blockseqrepo"
	mock_debondingdelegationseqrepo "github.com/figment-networks/oasishub-indexer/mock/repos/debondingdelegationseqrepo"
	mock_delegationseqrepo "github.com/figment-networks/oasishub-indexer/mock/repos/delegationseqrepo"
	mock_stakingseqrepo "github.com/figment-networks/oasishub-indexer/mock/repos/stakingseqrepo"
	mock_transactionseqrepo "github.com/figment-networks/oasishub-indexer/mock/repos/transactionseqrepo"
	mock_validatorseqrepo "github.com/figment-networks/oasishub-indexer/mock/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func Test_Sequencer_Block(t *testing.T) {
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

		blockDbRepoMock := mock_blockseqrepo.NewMockDbRepo(ctrl)
		blockDbRepoMock.EXPECT().GetByHeight(startH).Return(nil, errors.NewErrorFromMessage("not found", errors.NotFoundError)).Times(1)
		blockDbRepoMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

		validatorDbRepoMock := mock_validatorseqrepo.NewMockDbRepo(ctrl)
		validatorDbRepoMock.EXPECT().GetByHeight(startH).Times(1)
		validatorDbRepoMock.EXPECT().Create(gomock.Any()).Return(nil).MinTimes(1)

		stakingDbRepoMock := mock_stakingseqrepo.NewMockDbRepo(ctrl)
		stakingDbRepoMock.EXPECT().GetByHeight(startH).Return(nil, errors.NewErrorFromMessage("not found", errors.NotFoundError))
		stakingDbRepoMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

		transactionDbRepoMock := mock_transactionseqrepo.NewMockDbRepo(ctrl)
		transactionDbRepoMock.EXPECT().GetByHeight(startH).Times(1)
		transactionDbRepoMock.EXPECT().Create(gomock.Any()).Return(nil).MinTimes(1)

		delegatorDbRepoMock := mock_delegationseqrepo.NewMockDbRepo(ctrl)
		delegatorDbRepoMock.EXPECT().GetByHeight(startH).Times(1)
		delegatorDbRepoMock.EXPECT().Create(gomock.Any()).Return(nil).MinTimes(1)

		debondingDelegatorDbRepoMock := mock_debondingdelegationseqrepo.NewMockDbRepo(ctrl)
		debondingDelegatorDbRepoMock.EXPECT().GetByHeight(startH).Times(1)
		debondingDelegatorDbRepoMock.EXPECT().Create(gomock.Any()).Return(nil).MinTimes(0)

		sequencer := NewSequencer(
			blockDbRepoMock,
			validatorDbRepoMock,
			transactionDbRepoMock,
			stakingDbRepoMock,
			delegatorDbRepoMock,
			debondingDelegatorDbRepoMock,
		)

		returnedPayload, err := sequencer.Process(ctx, pld)
		if err != nil {
			t.Errorf("should not return error. Err: %+v", err)
			return
		}
		p := returnedPayload.(*payload)

		if p.BlockSequence == nil {
			t.Errorf("payload.BlockSequence should be set")
		}
		if p.StakingSequence == nil {
			t.Errorf("payload.StakingSequenceCreator should be set")
		}
		if p.ValidatorSequences == nil {
			t.Errorf("payload.ValidatorSequences should be set")
		}
		if p.TransactionSequences == nil {
			t.Errorf("payload.TransactionSequencesCreator should be set")
		}
		if p.DelegationSequences == nil {
			t.Errorf("payload.DelegationSequences should be set")
		}
	})
}
