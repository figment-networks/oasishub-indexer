package indexer

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
)

func TestBlockSeqCreator_Run(t *testing.T) {
	setup()

	var currHeight int64 = 20
	var pHeight int64 = 18
	var extHeight int64 = 20

	pTime := *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))
	extTime := *types.NewTimeFromTime(time.Date(1987, 12, 11, 14, 0, 0, 0, time.UTC))

	var pCount int64 = 145
	var extCount int64 = 100

	testPayload := func() *payload {
		return &payload{
			CurrentHeight: currHeight,
			Syncable: &model.Syncable{
				Height: pHeight,
				Time:   pTime,
			},
			ParsedBlock: ParsedBlockData{
				TransactionsCount: pCount,
			},
		}
	}

	tests := []struct {
		description           string
		payload               *payload
		dbErr                 error
		expectErr             error
		expectUpdatedBlockSeq *model.BlockSeq
		expectNewBlockSeq     *model.BlockSeq
	}{
		{
			description: "updates UpdatedBlockSequence when block found",
			payload:     testPayload(),
			dbErr:       nil,
			expectErr:   nil,
			expectUpdatedBlockSeq: &model.BlockSeq{
				Sequence: &model.Sequence{
					Height: extHeight,
					Time:   extTime,
				},
				TransactionsCount: pCount,
			},
			expectNewBlockSeq: nil,
		},
		{
			description:           "updates NewBlockSequence when block not found",
			payload:               testPayload(),
			dbErr:                 store.ErrNotFound,
			expectErr:             nil,
			expectUpdatedBlockSeq: nil,
			expectNewBlockSeq: &model.BlockSeq{
				Sequence: &model.Sequence{
					Height: pHeight,
					Time:   pTime,
				},
				TransactionsCount: pCount,
			},
		},
		{
			description:           "returns error if unexpected database error",
			payload:               testPayload(),
			dbErr:                 errTestDbFind,
			expectErr:             errTestDbFind,
			expectUpdatedBlockSeq: nil,
			expectNewBlockSeq:     nil,
		},
		{
			description: "returns error if block is invalid",
			payload: &payload{
				CurrentHeight: currHeight,
				Syncable: &model.Syncable{
					Height: pHeight,
					Time:   pTime,
				},
				ParsedBlock: ParsedBlockData{
					TransactionsCount: -200,
				},
			},
			dbErr:                 nil,
			expectErr:             errInvalidBlockSeq,
			expectUpdatedBlockSeq: nil,
			expectNewBlockSeq:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockDb := mock.NewMockBlockSeqCreatorTaskStore(ctrl)

			if tt.dbErr != nil {
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, tt.dbErr).Times(1)
			} else {
				existing := &model.BlockSeq{
					Sequence: &model.Sequence{
						Height: extHeight,
						Time:   extTime,
					},
					TransactionsCount: extCount,
				}
				mockDb.EXPECT().FindByHeight(currHeight).Return(existing, nil).Times(1)
			}

			task := NewBlockSeqCreatorTask(mockDb)

			if err := task.Run(ctx, tt.payload); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if !reflect.DeepEqual(tt.payload.NewBlockSequence, tt.expectNewBlockSeq) {
				t.Errorf("unexpected NewBlockSequence, want: %+v, got: %+v", tt.expectNewBlockSeq, tt.payload.NewBlockSequence)
				return
			}

			if !reflect.DeepEqual(tt.payload.UpdatedBlockSequence, tt.expectUpdatedBlockSeq) {
				t.Errorf("unexpected UpdatedBlockSequence, want: %+v, got: %+v", tt.expectUpdatedBlockSeq, tt.payload.UpdatedBlockSequence)
				return
			}
		})
	}
}

func TestValidatorSeqCreator_Run(t *testing.T) {
	setup()
	isTrue := true
	const currHeight int64 = 20
	const pHeight int64 = 18
	const extHeight int64 = 70

	pTime := *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))
	extTime := *types.NewTimeFromTime(time.Date(2018, 10, 10, 10, 0, 0, 0, time.UTC))

	newValidatorSeq := func(key string, addr string, power int64, height int64, _time types.Time) *model.ValidatorSeq {
		return &model.ValidatorSeq{
			Sequence: &model.Sequence{
				Height: height,
				Time:   _time,
			},
			EntityUID:   key,
			Address:     addr,
			VotingPower: power,
		}
	}

	tests := []struct {
		description string
		raw         []*validatorpb.Validator
		parsed      ParsedValidatorsData
		dbErr       error
		expectErr   error
	}{
		{
			description: "updates payload validator sequences",
			raw:         []*validatorpb.Validator{testpbValidator()},
			parsed:      make(ParsedValidatorsData),
			dbErr:       nil,
			expectErr:   nil,
		},
		{
			description: "return error if there's an unexpected database error",
			raw:         []*validatorpb.Validator{testpbValidator()},
			parsed:      make(ParsedValidatorsData),
			dbErr:       errTestDbFind,
			expectErr:   errTestDbFind,
		},
		{
			description: "updates validators with parsed data",
			raw: []*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
				testpbValidator(setValidatorAddress("addr2")),
			},
			parsed: ParsedValidatorsData{
				"addr0": parsedValidator{
					Proposed:           false,
					PrecommitValidated: &isTrue,
					TotalShares:        types.NewQuantity(big.NewInt(50)),
				},
				"addr1": parsedValidator{
					Proposed:           true,
					PrecommitValidated: &isTrue,
					TotalShares:        types.NewQuantity(big.NewInt(67)),
				},
			},
			dbErr:     nil,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[new sequences] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorSeqCreatorTaskStore(ctrl)

			expectNew := make([]*model.ValidatorSeq, len(tt.raw))
			for i, raw := range tt.raw {
				key := raw.GetNode().GetEntityId()
				if tt.expectErr == errTestDbFind {
					dbMock.EXPECT().FindByHeightAndEntityUID(currHeight, key).Return(nil, errTestDbFind).Times(1)
					break
				}
				dbMock.EXPECT().FindByHeightAndEntityUID(currHeight, key).Return(nil, store.ErrNotFound).Times(1)

				validator := newValidatorSeq(key, raw.GetAddress(), raw.GetVotingPower(), pHeight, pTime)
				if parsed, ok := tt.parsed[validator.Address]; ok {
					updateParsedValidatorSeq(validator, parsed)
				}
				expectNew[i] = validator
			}

			task := NewValidatorSeqCreatorTask(dbMock)

			pl := &payload{
				CurrentHeight: currHeight,
				Syncable: &model.Syncable{
					Height: pHeight,
					Time:   pTime,
				},
				RawValidators:    tt.raw,
				ParsedValidators: tt.parsed,
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if len(pl.NewValidatorSequences) != len(tt.raw) {
				t.Errorf("expected payload.NewValidatorSequences to contain new validators, got: %v; want: %v", len(pl.NewValidatorSequences), len(tt.raw))
				return
			}

			for _, expectVal := range expectNew {
				var found bool
				for _, val := range pl.NewValidatorSequences {
					if val.Address == expectVal.Address {
						if !reflect.DeepEqual(val, *expectVal) {
							t.Errorf("unexpected entry in payload.NewValidatorSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.NewValidatorSequences, want: %v", expectVal)
				}
			}
		})

		t.Run(fmt.Sprintf("[old sequences] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorSeqCreatorTaskStore(ctrl)

			expectUpdated := make([]*model.ValidatorSeq, len(tt.raw))
			for i, raw := range tt.raw {
				key := raw.GetNode().GetEntityId()
				if tt.dbErr != nil {
					dbMock.EXPECT().FindByHeightAndEntityUID(currHeight, key).Return(nil, tt.dbErr).Times(1)
					break
				}
				existing := newValidatorSeq(key, "existingAddr", 1000, extHeight, extTime)
				dbMock.EXPECT().FindByHeightAndEntityUID(currHeight, key).Return(existing, nil).Times(1)

				validator := newValidatorSeq(key, raw.GetAddress(), raw.GetVotingPower(), extHeight, extTime)
				if parsed, ok := tt.parsed[validator.Address]; ok {
					updateParsedValidatorSeq(validator, parsed)
				}
				expectUpdated[i] = validator
			}

			task := NewValidatorSeqCreatorTask(dbMock)

			pl := &payload{
				CurrentHeight: currHeight,
				Syncable: &model.Syncable{
					Height: pHeight,
					Time:   pTime,
				},
				RawValidators:    tt.raw,
				ParsedValidators: tt.parsed,
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if len(pl.UpdatedValidatorSequences) != len(tt.raw) {
				t.Errorf("expected payload.UpdatedValidatorSequences to contain new validators, got: %v; want: %v", len(pl.UpdatedValidatorSequences), len(tt.raw))
				return
			}

			for _, expectVal := range expectUpdated {
				var found bool
				for _, val := range pl.UpdatedValidatorSequences {
					if val.Address == expectVal.Address {
						if !reflect.DeepEqual(val, *expectVal) {
							t.Errorf("unexpected entry in payload.UpdatedValidatorSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.UpdatedValidatorSequences, want: %v", expectVal)
				}
			}
		})
	}
}

func TestStakingSeqCreator_Run(t *testing.T) {
	setup()
	var currHeight int64 = 20

	sync := &model.Syncable{
		Height: 18,
		Time:   *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC)),
	}

	existing := &model.StakingSeq{
		Sequence: &model.Sequence{
			Height: 20,
			Time:   *types.NewTimeFromTime(time.Date(1987, 12, 11, 14, 0, 0, 0, time.UTC)),
		},
		TotalSupply:         types.NewQuantityFromBytes(randBytes(6)),
		CommonPool:          types.NewQuantityFromBytes(randBytes(6)),
		DebondingInterval:   rand.Uint64(),
		MinDelegationAmount: types.NewQuantityFromBytes(randBytes(6)),
	}

	tests := []struct {
		description      string
		raw              *statepb.Staking
		dbErr            error
		expectErr        error
		expectStakingSeq *model.StakingSeq
	}{
		{
			description:      "Adds exisitng staking seq to payload",
			raw:              testpbStaking(),
			dbErr:            nil,
			expectErr:        nil,
			expectStakingSeq: existing,
		},
		{
			description: "Adds new staking seq to payload",
			raw: testpbStaking(
				setStakingTotalSupply([]byte{1}),
				setStakingCommonPool([]byte{2}),
				setStakingDebondingInterval(3),
				setStakingMinDelegationAmount([]byte{4}),
			),
			dbErr:     store.ErrNotFound,
			expectErr: nil,
			expectStakingSeq: &model.StakingSeq{
				Sequence: &model.Sequence{
					Height: sync.Height,
					Time:   sync.Time,
				},
				TotalSupply:         types.NewQuantityFromBytes([]byte{1}),
				CommonPool:          types.NewQuantityFromBytes([]byte{2}),
				DebondingInterval:   3,
				MinDelegationAmount: types.NewQuantityFromBytes([]byte{4}),
			},
		},
		{
			description:      "Returns error on unexpected FindByHeight database error",
			raw:              testpbStaking(),
			dbErr:            errTestDbFind,
			expectErr:        errTestDbFind,
			expectStakingSeq: nil,
		},
		{
			description: "Returns error on unexpected Create database error",
			raw: testpbStaking(
				setStakingTotalSupply([]byte{1}),
				setStakingCommonPool([]byte{2}),
				setStakingDebondingInterval(3),
				setStakingMinDelegationAmount([]byte{4}),
			),
			dbErr:     errTestDbCreate,
			expectErr: errTestDbCreate,
			expectStakingSeq: &model.StakingSeq{
				Sequence: &model.Sequence{
					Height: sync.Height,
					Time:   sync.Time,
				},
				TotalSupply:         types.NewQuantityFromBytes([]byte{1}),
				CommonPool:          types.NewQuantityFromBytes([]byte{2}),
				DebondingInterval:   3,
				MinDelegationAmount: types.NewQuantityFromBytes([]byte{4}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockDb := mock.NewMockStakingSeqCreatorTaskStore(ctrl)

			if tt.dbErr == errTestDbFind {
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, errTestDbFind).Times(1)
			} else if tt.dbErr == errTestDbCreate {
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, store.ErrNotFound).Times(1)
				mockDb.EXPECT().Create(tt.expectStakingSeq).Return(errTestDbCreate).Times(1)
			} else if tt.dbErr == store.ErrNotFound {
				// expet new seq to be created and added to payload
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, store.ErrNotFound).Times(1)
				mockDb.EXPECT().Create(tt.expectStakingSeq).Return(nil).Times(1)
			} else {
				// expect existing seq to be added to payload
				mockDb.EXPECT().FindByHeight(currHeight).Return(existing, nil).Times(1)
			}

			task := NewStakingSeqCreatorTask(mockDb)
			pl := &payload{
				CurrentHeight: currHeight,
				Syncable:      sync,
				RawState: &statepb.State{
					Staking: tt.raw,
				},
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if !reflect.DeepEqual(pl.StakingSequence, tt.expectStakingSeq) {
				t.Errorf("unexpected NewBlockSequence, want: %+v, got: %+v", tt.expectStakingSeq, pl.StakingSequence)
				return
			}
		})
	}
}

func TestTransactionSeqCreator_Run(t *testing.T) {
	setup()
	var currHeight int64 = 18

	sync := &model.Syncable{
		Height: 18,
		Time:   *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC)),
	}
	seq := &model.Sequence{
		Height: sync.Height,
		Time:   sync.Time,
	}

	emptyRaw := make([]*transactionpb.Transaction, 0)
	// emptyModel := make([]model.TransactionSeq, 0)

	rawToModel := func(raw *transactionpb.Transaction) *model.TransactionSeq {
		return &model.TransactionSeq{
			Sequence:  seq,
			PublicKey: raw.GetPublicKey(),
			Hash:      raw.GetHash(),
			GasPrice:  types.NewQuantityFromBytes(raw.GetGasPrice()),
		}
	}

	raw1 := testpbTransaction("raw1")
	raw2 := testpbTransaction("raw2")
	raw3 := testpbTransaction("raw3")

	tests := []struct {
		description string
		rawExisting []*transactionpb.Transaction
		rawNew      []*transactionpb.Transaction

		dbErr     error
		expectErr error
		expectSeq []*model.TransactionSeq
	}{
		{
			description: "Adds exisitng transaction seq to payload",
			rawExisting: []*transactionpb.Transaction{raw1},
			rawNew:      emptyRaw,
			dbErr:       nil,
			expectErr:   nil,
			expectSeq:   []*model.TransactionSeq{rawToModel(raw1)},
		},
		{
			description: "Adds new transaction seq to payload",
			rawExisting: emptyRaw,
			rawNew:      []*transactionpb.Transaction{raw1},
			dbErr:       nil,
			expectErr:   nil,
			expectSeq:   []*model.TransactionSeq{rawToModel(raw1)},
		},
		{
			description: "Returns err on unexpected FindByHeight error",
			rawExisting: emptyRaw,
			rawNew:      []*transactionpb.Transaction{raw1},
			dbErr:       errTestDbFind,
			expectErr:   errTestDbFind,
			expectSeq:   []*model.TransactionSeq{},
		},
		{
			description: "Adds empty list to payload when there's no transaction seq",
			rawExisting: emptyRaw,
			rawNew:      emptyRaw,
			dbErr:       nil,
			expectErr:   nil,
			expectSeq:   []*model.TransactionSeq{},
		},
		{
			description: "Adds new and existing transaction sequences to payload",
			rawExisting: []*transactionpb.Transaction{raw1, raw2},
			rawNew:      []*transactionpb.Transaction{raw3},
			dbErr:       nil,
			expectErr:   nil,
			expectSeq:   []*model.TransactionSeq{rawToModel(raw1), rawToModel(raw2), rawToModel(raw3)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			mockDb := mock.NewMockTransactionSeqCreatorTaskStore(ctrl)

			if tt.dbErr != nil {
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, tt.dbErr).Times(1)
			} else {
				dbReturn := make([]model.TransactionSeq, len(tt.rawExisting))
				for i, raw := range tt.rawExisting {
					dbReturn[i] = *rawToModel(raw)
				}

				mockDb.EXPECT().FindByHeight(currHeight).Return(dbReturn, nil).Times(1)

				// expect create to be called on new raw transaction seq
				for _, raw := range tt.rawNew {
					mockDb.EXPECT().Create(rawToModel(raw)).Return(nil).Times(1)
				}
			}

			task := NewTransactionSeqCreatorTask(mockDb)
			pl := &payload{
				CurrentHeight:   currHeight,
				Syncable:        sync,
				RawTransactions: append(tt.rawExisting, tt.rawNew...),
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if len(pl.TransactionSequences) != len(tt.expectSeq) {
				t.Errorf("expected payload.TransactionSequences to contain all sequences, got: %v; want: %v", len(pl.TransactionSequences), len(tt.expectSeq))
				return
			}

			for _, expectVal := range tt.expectSeq {
				var found bool
				for _, val := range pl.TransactionSequences {
					if val.PublicKey == expectVal.PublicKey {
						if !reflect.DeepEqual(val, *expectVal) {
							t.Errorf("unexpected entry in payload.TransactionSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.TransactionSequences, want: %v", expectVal)
				}
			}
		})
	}
}

func TestDelegationSeqCreator_Run(t *testing.T) {
	setup()
	var currHeight int64 = 18

	sync := &model.Syncable{
		Height: currHeight,
		Time:   *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC)),
	}

	toModel := func(valUID, delUID string, shares []byte) model.DelegationSeq {
		return model.DelegationSeq{
			Sequence: &model.Sequence{
				Height: sync.Height,
				Time:   sync.Time,
			},

			ValidatorUID: valUID,
			DelegatorUID: delUID,
			Shares:       types.NewQuantityFromBytes(shares),
		}
	}

	tests := []struct {
		description string
		rawStaking  *statepb.Staking
		dbReturn    []model.DelegationSeq
		dbErr       error
		expectErr   error
		expectSeq   []model.DelegationSeq
	}{
		{
			description: "Adds exisitng delegation seq to payload",
			rawStaking: testpbStaking(
				setStakingDelegationEntry("t0", "del1", uintToBytes(100, t)),
			),
			dbReturn:  []model.DelegationSeq{toModel("t0", "del1", uintToBytes(100, t))},
			dbErr:     nil,
			expectErr: nil,
			expectSeq: []model.DelegationSeq{toModel("t0", "del1", uintToBytes(100, t))},
		},
		{
			description: "Adds new delegation seq to payload",
			rawStaking: testpbStaking(
				setStakingDelegationEntry("t0", "del1", uintToBytes(100, t)),
			),
			dbReturn:  []model.DelegationSeq{},
			dbErr:     nil,
			expectErr: nil,
			expectSeq: []model.DelegationSeq{toModel("t0", "del1", uintToBytes(100, t))},
		},
		{
			description: "Returns err on unexpected FindByHeight error",
			rawStaking: testpbStaking(
				setStakingDelegationEntry("t0", "del1", uintToBytes(100, t)),
			),
			dbReturn:  []model.DelegationSeq{},
			dbErr:     errTestDbFind,
			expectErr: errTestDbFind,
			expectSeq: []model.DelegationSeq{},
		},
		{
			description: "Adds empty list to payload when there's no delegations sequence",
			rawStaking:  testpbStaking(),
			dbReturn:    []model.DelegationSeq{},
			dbErr:       nil,
			expectErr:   nil,
			expectSeq:   []model.DelegationSeq{},
		},
		{
			description: "Adds new and existing delegations sequence to payload",
			rawStaking: testpbStaking(
				setStakingDelegationEntry("t0", "del1", uintToBytes(100, t)),
				setStakingDelegationEntry("t0", "newdel", uintToBytes(200, t)),
				setStakingDelegationEntry("t3", "newdel2", uintToBytes(300, t)),
				setStakingDelegationEntry("t1", "del1", uintToBytes(400, t)),
			),
			dbReturn: []model.DelegationSeq{
				toModel("t0", "del1", uintToBytes(100, t)),
				toModel("t1", "del1", uintToBytes(400, t)),
			},
			dbErr:     nil,
			expectErr: nil,
			expectSeq: []model.DelegationSeq{
				toModel("t0", "del1", uintToBytes(100, t)),
				toModel("t0", "newdel", uintToBytes(200, t)),
				toModel("t3", "newdel2", uintToBytes(300, t)),
				toModel("t1", "del1", uintToBytes(400, t)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			mockDb := mock.NewMockDelegationSeqCreatorTaskStore(ctrl)

			if tt.dbErr != nil {
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, tt.dbErr).Times(1)
			} else {
				mockDb.EXPECT().FindByHeight(currHeight).Return(tt.dbReturn, nil).Times(1)

				// expect create to be called on each delegation seq not in dbReturn
				for _, seq := range tt.expectSeq {
					var found bool
					for _, extSeq := range tt.dbReturn {
						if reflect.DeepEqual(seq, extSeq) {
							found = true
							break
						}
					}
					if !found {
						s := seq
						mockDb.EXPECT().Create(&s).Return(nil).Times(1)
					}
				}
			}

			task := NewDelegationsSeqCreatorTask(mockDb)
			pl := &payload{
				CurrentHeight: currHeight,
				Syncable:      sync,
				RawState: &statepb.State{
					Staking: tt.rawStaking,
				},
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if len(pl.DelegationSequences) != len(tt.expectSeq) {
				t.Errorf("expected payload.DelegationSequences to contain all sequences, got: %v; want: %v", len(pl.DelegationSequences), len(tt.expectSeq))
				return
			}

			for _, expectVal := range tt.expectSeq {
				var found bool
				for _, val := range pl.DelegationSequences {
					if val.DelegatorUID == expectVal.DelegatorUID && val.ValidatorUID == expectVal.ValidatorUID {
						if !reflect.DeepEqual(val, expectVal) {
							t.Errorf("unexpected entry in payload.DelegationSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.DelegationSequences, want: %v", expectVal)
				}
			}
		})
	}
}

func TestDebondingDelegationSeqCreator_Run(t *testing.T) {
	setup()
	var currHeight int64 = 18

	sync := &model.Syncable{
		Height: currHeight,
		Time:   *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC)),
	}

	toModel := func(valUID, delUID string, shares []byte, endTime uint64) model.DebondingDelegationSeq {
		return model.DebondingDelegationSeq{
			Sequence: &model.Sequence{
				Height: sync.Height,
				Time:   sync.Time,
			},

			ValidatorUID: valUID,
			DelegatorUID: delUID,
			Shares:       types.NewQuantityFromBytes(shares),
			DebondEnd:    endTime,
		}
	}

	tests := []struct {
		description string
		rawStaking  *statepb.Staking
		dbReturn    []model.DebondingDelegationSeq
		dbErr       error
		expectErr   error
		expectSeq   []model.DebondingDelegationSeq
	}{
		{
			description: "Adds exisitng delegation seq to payload",
			rawStaking: testpbStaking(
				setDebondingDelegationEntry("t0", "del1", uintToBytes(100, t), 12),
			),
			dbReturn:  []model.DebondingDelegationSeq{toModel("t0", "del1", uintToBytes(100, t), 12)},
			dbErr:     nil,
			expectErr: nil,
			expectSeq: []model.DebondingDelegationSeq{toModel("t0", "del1", uintToBytes(100, t), 12)},
		},
		{
			description: "Adds new delegation seq to payload",
			rawStaking: testpbStaking(
				setDebondingDelegationEntry("t0", "del1", uintToBytes(100, t), 14),
			),
			dbReturn:  []model.DebondingDelegationSeq{},
			dbErr:     nil,
			expectErr: nil,
			expectSeq: []model.DebondingDelegationSeq{toModel("t0", "del1", uintToBytes(100, t), 14)},
		},
		{
			description: "Returns err on unexpected FindByHeight error",
			rawStaking: testpbStaking(
				setDebondingDelegationEntry("t0", "del1", uintToBytes(100, t), 5),
			),
			dbReturn:  []model.DebondingDelegationSeq{},
			dbErr:     errTestDbFind,
			expectErr: errTestDbFind,
			expectSeq: []model.DebondingDelegationSeq{},
		},
		{
			description: "Adds empty list to payload when there's no delegations sequence",
			rawStaking:  testpbStaking(),
			dbReturn:    []model.DebondingDelegationSeq{},
			dbErr:       nil,
			expectErr:   nil,
			expectSeq:   []model.DebondingDelegationSeq{},
		},
		{
			description: "Adds new and existing delegations sequence to payload",
			rawStaking: testpbStaking(
				setDebondingDelegationEntry("t0", "del1", uintToBytes(100, t), 1),
				setDebondingDelegationEntry("t0", "newdel", uintToBytes(200, t), 2),
				setDebondingDelegationEntry("t3", "newdel2", uintToBytes(300, t), 3),
				setDebondingDelegationEntry("t1", "del1", uintToBytes(400, t), 4),
			),
			dbReturn: []model.DebondingDelegationSeq{
				toModel("t0", "del1", uintToBytes(100, t), 1),
				toModel("t1", "del1", uintToBytes(400, t), 4),
			},
			dbErr:     nil,
			expectErr: nil,
			expectSeq: []model.DebondingDelegationSeq{
				toModel("t0", "del1", uintToBytes(100, t), 1),
				toModel("t0", "newdel", uintToBytes(200, t), 2),
				toModel("t3", "newdel2", uintToBytes(300, t), 3),
				toModel("t1", "del1", uintToBytes(400, t), 4),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			mockDb := mock.NewMockDebondingDelegationSeqCreatorTaskStore(ctrl)

			if tt.dbErr != nil {
				mockDb.EXPECT().FindByHeight(currHeight).Return(nil, tt.dbErr).Times(1)
			} else {
				mockDb.EXPECT().FindByHeight(currHeight).Return(tt.dbReturn, nil).Times(1)

				// expect create to be called on each delegation seq not in dbReturn
				for _, seq := range tt.expectSeq {
					var found bool
					for _, extSeq := range tt.dbReturn {
						if reflect.DeepEqual(seq, extSeq) {
							found = true
							break
						}
					}
					if !found {
						s := seq
						mockDb.EXPECT().Create(&s).Return(nil).Times(1)
					}
				}
			}

			task := NewDebondingDelegationsSeqCreatorTask(mockDb)
			pl := &payload{
				CurrentHeight: currHeight,
				Syncable:      sync,
				RawState: &statepb.State{
					Staking: tt.rawStaking,
				},
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr != nil {
				return
			}

			if len(pl.DebondingDelegationSequences) != len(tt.expectSeq) {
				t.Errorf("expected payload.DebondingDelegationSequences to contain all sequences, got: %v; want: %v", len(pl.DelegationSequences), len(tt.expectSeq))
				return
			}

			for _, expectVal := range tt.expectSeq {
				var found bool
				for _, val := range pl.DebondingDelegationSequences {
					if val.DelegatorUID == expectVal.DelegatorUID && val.ValidatorUID == expectVal.ValidatorUID {
						if !reflect.DeepEqual(val, expectVal) {
							t.Errorf("unexpected entry in payload.DebondingDelegationSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.DebondingDelegationSequences, want: %v", expectVal)
				}
			}
		})
	}
}

func updateParsedValidatorSeq(m *model.ValidatorSeq, parsed parsedValidator) {
	m.PrecommitValidated = parsed.PrecommitValidated
	m.Proposed = parsed.Proposed
	m.TotalShares = parsed.TotalShares
}
