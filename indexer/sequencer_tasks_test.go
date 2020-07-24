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
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
)

func TestBlockSeqCreator_Run(t *testing.T) {
	setup(t)

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
	setup(t)
	isTrue := true
	const currHeight int64 = 20
	const pHeight int64 = 18
	const extHeight int64 = 70

	pTime := *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))
	extTime := *types.NewTimeFromTime(time.Date(2018, 10, 10, 10, 0, 0, 0, time.UTC))

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
	setup(t)
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

func newValidatorSeq(key string, addr string, power int64, height int64, _time types.Time) *model.ValidatorSeq {
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

func updateParsedValidatorSeq(m *model.ValidatorSeq, parsed parsedValidator) {
	m.PrecommitValidated = parsed.PrecommitValidated
	m.Proposed = parsed.Proposed
	m.TotalShares = parsed.TotalShares
}
