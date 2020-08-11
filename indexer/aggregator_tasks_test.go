package indexer

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
)

type accountLedger map[string]*accountpb.Account

func TestAccountAggCreatorTask_Run(t *testing.T) {
	setup()

	tests := []struct {
		description string
		new         accountLedger
		existing    accountLedger
		expectErr   error
	}{
		{
			description: "creates new accounts",
			new: accountLedger{
				"pkey1": testAccount(),
				"pkey2": testAccount(),
			},
			existing:  accountLedger{},
			expectErr: nil,
		},
		{
			description: "updates existing accounts",
			new:         accountLedger{},
			existing: accountLedger{
				"pkey1": testAccount(),
				"pkey2": testAccount(),
			},
			expectErr: nil,
		},
		{
			description: "creates and updates accounts",
			new: accountLedger{
				"pkey3": testAccount(),
			},
			existing: accountLedger{
				"pkey1": testAccount(),
				"pkey2": testAccount(),
			},
			expectErr: nil,
		},
		{
			description: "return error if there's an unexpected db error on findByPublicKey",
			new: accountLedger{
				"pkey1": testAccount(),
			},
			existing:  accountLedger{},
			expectErr: errTestDbFind,
		},
		{
			description: "return error if there's a db error on create",
			new: accountLedger{
				"pkey1": testAccount(),
			},
			existing:  accountLedger{},
			expectErr: errTestDbCreate,
		},
		{
			description: "return error if there's a db error on save",
			new:         accountLedger{},
			existing: accountLedger{
				"pkey1": testAccount(),
			},
			expectErr: errTestDbSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockAccountAggCreatorTaskStore(ctrl)

			ledger := combineLedgers(tt.new, tt.existing)
			payload := testAccountAggPayload(ledger)

			for key, acnt := range tt.new {
				if tt.expectErr == errTestDbFind {
					dbMock.EXPECT().FindByPublicKey(key).Return(nil, errTestDbFind).Times(1)
					break
				}
				dbMock.EXPECT().FindByPublicKey(key).Return(nil, store.ErrNotFound).Times(1)
				newAccount := newAccountAgg(key, payload.Syncable.Height, payload.Syncable.Time)
				updatedAccount := updateAccountAgg(newAccount, acnt, payload)
				if tt.expectErr == errTestDbCreate {
					dbMock.EXPECT().Create(updatedAccount).Return(errTestDbCreate).Times(1)
					break
				}
				dbMock.EXPECT().Create(updatedAccount).Return(nil).Times(1)
			}

			for key, acnt := range tt.existing {
				existingAccount := newAccountAgg(key, 0, *types.NewTimeFromTime(time.Now()))
				dbMock.EXPECT().FindByPublicKey(key).Return(existingAccount, nil).Times(1)
				updatedAcnt := updateAccountAgg(existingAccount, acnt, payload)

				if tt.expectErr == errTestDbSave {
					dbMock.EXPECT().Save(updatedAcnt).Return(errTestDbSave).Times(1)
					break
				}

				dbMock.EXPECT().Save(updatedAcnt).Return(nil).Times(1)
			}

			task := NewAccountAggCreatorTask(dbMock)
			if err := task.Run(ctx, payload); err != tt.expectErr {
				t.Errorf("unexpected error, got: %v; want: %v", err, tt.expectErr)
				return
			}

			// don't check payload if expected error
			if tt.expectErr != nil {
				return
			}

			if len(payload.NewAggregatedAccounts) != len(tt.new) {
				t.Errorf("expected payload.NewAggregatedAccounts to contain new accounts, got: %v; want: %v", len(payload.NewAggregatedAccounts), len(tt.new))
				return
			}

			if len(payload.UpdatedAggregatedAccounts) != len(tt.existing) {
				t.Errorf("expected payload.UpdatedAggregatedAccounts to contain accounts, got: %v; want: %v", len(payload.UpdatedAggregatedAccounts), len(tt.existing))
				return
			}
		})
	}
}

func TestValidatorAggCreatorTask_Run(t *testing.T) {
	setup()
	plTime := *types.NewTimeFromTime(time.Now())
	const syncHeight int64 = 17
	const currHeight int64 = 64

	const nilValidated int64 = 0
	const notValidated int64 = 1
	const validated int64 = 2

	tests := []struct {
		description string
		raw         []*validatorpb.Validator
		parsed      ParsedValidatorsData
		expectErr   error
	}{
		{
			"updates payload with accounts",
			[]*validatorpb.Validator{testpbValidator()},
			make(ParsedValidatorsData),
			nil,
		},
		{
			"return error if there's an unexpected db error on errTestDbFind",
			[]*validatorpb.Validator{testpbValidator(), testpbValidator()},
			make(ParsedValidatorsData),
			errTestDbFind,
		},
		{
			"updates validator with nil validated data",
			[]*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
			},
			ParsedValidatorsData{
				"addr1": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: nilValidated,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
			},
			nil,
		},
		{
			"updates validator with not validated data",
			[]*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
			},
			ParsedValidatorsData{
				"addr1": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: notValidated,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
			},
			nil,
		},
		{
			"updates validator with validated data",
			[]*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
			},
			ParsedValidatorsData{
				"addr1": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: validated,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
			},
			nil,
		},
		{
			"updates validator with proposed data",
			[]*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
			},
			ParsedValidatorsData{
				"addr1": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: validated,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
			},
			nil,
		},
		{
			"updates multiple validators",
			[]*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
				testpbValidator(setValidatorAddress("addr2")),
				testpbValidator(setValidatorAddress("addr3")),
			},
			ParsedValidatorsData{
				"addr1": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: 1,
					TotalShares:          types.NewQuantity(big.NewInt(66)),
				},
				"addr2": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 2,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
				"addr3": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 0,
					TotalShares:          types.NewQuantity(big.NewInt(68)),
				},
				"bonus_addr": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: 0,
					TotalShares:          types.NewQuantity(big.NewInt(68)),
				},
			},
			nil,
		},
		{
			"updates validators with parsedValidator data",
			[]*validatorpb.Validator{
				testpbValidator(setValidatorAddress("addr1")),
				testpbValidator(setValidatorAddress("addr2")),
				testpbValidator(setValidatorAddress("addr3")),
			},
			ParsedValidatorsData{
				"addr1": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: 1,
					TotalShares:          types.NewQuantity(big.NewInt(66)),
				},
				"addr2": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 2,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
				"addr3": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 0,
					TotalShares:          types.NewQuantity(big.NewInt(68)),
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[new validators] %s", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorAggCreatorTaskStore(ctrl)

			payload := &payload{
				Syncable: &model.Syncable{
					Height: syncHeight,
					Time:   plTime,
				},
				CurrentHeight:    currHeight,
				RawValidators:    tt.raw,
				ParsedValidators: tt.parsed,
			}

			expectValidators := make([]*model.ValidatorAgg, len(tt.raw))
			for i, raw := range tt.raw {
				key := raw.GetNode().GetEntityId()
				if tt.expectErr == errTestDbFind {
					dbMock.EXPECT().FindByEntityUID(key).Return(nil, errTestDbFind).Times(1)
					break
				}
				dbMock.EXPECT().FindByEntityUID(key).Return(nil, store.ErrNotFound).Times(1)

				validator := newValidatorAgg(key, raw.GetAddress(), payload.Syncable.Height, payload.Syncable.Time)
				validator = updateValidatorAgg(validator, raw, payload)

				if parsed, ok := tt.parsed[validator.Address]; ok {
					updateParsedValidatorAgg(validator, parsed, payload, true)
				}

				expectValidators[i] = validator
			}

			task := NewValidatorAggCreatorTask(dbMock)
			if err := task.Run(ctx, payload); err != tt.expectErr {
				t.Errorf("unexpected error, got: %v; want: %v", err, tt.expectErr)
				return
			}

			// don't check payload if there was an error
			if tt.expectErr != nil {
				return
			}

			if len(payload.NewAggregatedValidators) != len(tt.raw) {
				t.Errorf("expected payload.NewAggregatedValidators to contain new validators, got: %v; want: %v", len(payload.NewAggregatedValidators), len(tt.raw))
				return
			}

			for _, expectVal := range expectValidators {
				var found bool
				for _, val := range payload.NewAggregatedValidators {
					if val.Address == expectVal.Address {
						if !reflect.DeepEqual(val, *expectVal) {
							t.Errorf("unexpected entry in payload.NewAggregatedValidators, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.NewAggregatedValidators, want: %v", expectVal)
				}
			}
		})
	}

	for _, tt := range tests {
		t.Run((fmt.Sprintf("[existing validators] %s", tt.description)), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorAggCreatorTaskStore(ctrl)

			payload := &payload{
				Syncable: &model.Syncable{
					Height: syncHeight,
					Time:   plTime,
				},
				CurrentHeight:    currHeight,
				RawValidators:    tt.raw,
				ParsedValidators: tt.parsed,
			}

			expectExisting := make([]*model.ValidatorAgg, len(tt.raw))
			for i, raw := range tt.raw {
				key := raw.GetNode().GetEntityId()
				if tt.expectErr == errTestDbFind {
					dbMock.EXPECT().FindByEntityUID(key).Return(nil, errTestDbFind).Times(1)
					break
				}

				returnVal := newValidatorAgg(key, randString(5), 0, *types.NewTimeFromTime(time.Now()))
				dbMock.EXPECT().FindByEntityUID(key).Return(returnVal, nil).Times(1)

				validator := updateValidatorAgg(returnVal, raw, payload)
				validator.RecentTendermintAddress = raw.GetAddress()

				if parsed, ok := tt.parsed[raw.Address]; ok {
					updateParsedValidatorAgg(validator, parsed, payload, false)
				}
				expectExisting[i] = validator
			}

			task := NewValidatorAggCreatorTask(dbMock)
			if err := task.Run(ctx, payload); err != tt.expectErr {
				t.Errorf("unexpected error, got: %v; want: %v", err, tt.expectErr)
				return
			}

			// don't check payload if there was an error
			if tt.expectErr != nil {
				return
			}

			if len(payload.UpdatedAggregatedValidators) != len(tt.raw) {
				t.Errorf("expected payload.UpdatedAggregatedValidators to contain accounts, got: %v; want: %v", len(payload.UpdatedAggregatedValidators), len(tt.raw))
				return
			}

			for _, expectVal := range expectExisting {
				var found bool
				for _, val := range payload.UpdatedAggregatedValidators {
					if val.Address == expectVal.Address {
						if !reflect.DeepEqual(val, *expectVal) {
							t.Errorf("unexpected entry in payload.UpdatedAggregatedValidators, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.UpdatedAggregatedValidators, want: %v", expectVal)
				}
			}
		})
	}
}

func testAccountAggPayload(ledger accountLedger) *payload {
	return &payload{
		CurrentHeight: 10,
		Syncable: &model.Syncable{
			Height: 10,
			Time:   *types.NewTimeFromTime(time.Now()),
		},
		RawState: &statepb.State{
			Staking: &statepb.Staking{
				Ledger: ledger,
			},
		},
	}
}

func testAccount() *accountpb.Account {
	return &accountpb.Account{
		General: &accountpb.GeneralAccount{
			Balance: randBytes(10),
			Nonce:   100,
		},
		Escrow: &accountpb.EscrowAccount{
			Active: &accountpb.SharePool{
				Balance:     randBytes(10),
				TotalShares: randBytes(10),
			},
			Debonding: &accountpb.SharePool{
				Balance:     randBytes(10),
				TotalShares: randBytes(10),
			},
		},
	}
}

func newAccountAgg(key string, height int64, _time types.Time) *model.AccountAgg {
	return &model.AccountAgg{
		Aggregate: &model.Aggregate{
			StartedAtHeight: height,
			StartedAt:       _time,
		},
		PublicKey: key,
	}
}

func updateAccountAgg(original *model.AccountAgg, acnt *accountpb.Account, pl *payload) *model.AccountAgg {
	m := &model.AccountAgg{
		Aggregate: &model.Aggregate{
			StartedAtHeight: original.Aggregate.StartedAtHeight,
			StartedAt:       original.Aggregate.StartedAt,
			RecentAtHeight:  pl.Syncable.Height,
			RecentAt:        pl.Syncable.Time,
		},
		PublicKey: original.PublicKey,

		RecentGeneralBalance:             types.NewQuantityFromBytes(acnt.GetGeneral().GetBalance()),
		RecentGeneralNonce:               acnt.GetGeneral().GetNonce(),
		RecentEscrowActiveBalance:        types.NewQuantityFromBytes(acnt.GetEscrow().GetActive().GetBalance()),
		RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(acnt.GetEscrow().GetActive().GetTotalShares()),
		RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(acnt.GetEscrow().GetDebonding().GetBalance()),
		RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(acnt.GetEscrow().GetDebonding().GetTotalShares()),
	}
	return m
}

func newValidatorAgg(key string, addr string, height int64, _time types.Time) *model.ValidatorAgg {
	return &model.ValidatorAgg{
		Aggregate: &model.Aggregate{
			StartedAtHeight: height,
			StartedAt:       _time,
		},
		EntityUID: key,
		Address:   addr,
	}
}

func updateValidatorAgg(original *model.ValidatorAgg, raw *validatorpb.Validator, pl *payload) *model.ValidatorAgg {
	return &model.ValidatorAgg{
		Aggregate: &model.Aggregate{
			StartedAtHeight: original.Aggregate.StartedAtHeight,
			StartedAt:       original.Aggregate.StartedAt,
			RecentAtHeight:  pl.Syncable.Height,
			RecentAt:        pl.Syncable.Time,
		},
		EntityUID:               original.EntityUID,
		Address:                 original.Address,
		RecentTendermintAddress: raw.GetTendermintAddress(),
		RecentVotingPower:       raw.GetVotingPower(),
		RecentAsValidatorHeight: pl.Syncable.Height,
	}
}

func updateParsedValidatorAgg(m *model.ValidatorAgg, parsed parsedValidator, pl *payload, newValidator bool) {
	const notValidated int64 = 1
	const validated int64 = 2

	m.RecentTotalShares = parsed.TotalShares

	if parsed.PrecommitBlockIdFlag == notValidated {
		m.AccumulatedUptimeCount++
	} else if parsed.PrecommitBlockIdFlag == validated {
		m.AccumulatedUptime++
		m.AccumulatedUptimeCount++
	}

	if parsed.Proposed {
		if newValidator {
			m.RecentProposedHeight = pl.CurrentHeight
		} else {
			m.RecentProposedHeight = pl.Syncable.Height
		}
		m.AccumulatedProposedCount++
	}
}

func combineLedgers(m1, m2 accountLedger) accountLedger {
	ledger := make(accountLedger)
	for k, v := range m1 {
		ledger[k] = v
	}
	for k, v := range m2 {
		ledger[k] = v
	}
	return ledger
}
