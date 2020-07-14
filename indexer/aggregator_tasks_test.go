package indexer

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	mock "github.com/figment-networks/oasishub-indexer/indexer/mock"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

type accountLedger map[string]*accountpb.Account

func TestAggCreatorTask_Run(t *testing.T) {
	setup(t)

	findByPublicKeyErr := errors.New("findByPublicKeyErr")
	createErr := errors.New("createErr")
	saveErr := errors.New("saveErr")

	tests := []struct {
		description string
		new         accountLedger
		existing    accountLedger
		result      error
	}{
		{
			"creates new accounts",
			accountLedger{
				"pkey1": testAccount(),
				"pkey2": testAccount(),
			},
			accountLedger{},
			nil,
		},
		{
			"updates existing accounts",
			accountLedger{},
			accountLedger{
				"pkey1": testAccount(),
				"pkey2": testAccount(),
			},
			nil,
		},
		{
			"creates and updates accounts",
			accountLedger{
				"pkey3": testAccount(),
			},
			accountLedger{
				"pkey1": testAccount(),
				"pkey2": testAccount(),
			},
			nil,
		},
		{
			"return error if there's an unexpected db error on findByPublicKey",
			accountLedger{
				"pkey1": testAccount(),
			},
			accountLedger{},
			findByPublicKeyErr,
		},
		{
			"return error if there's a db error on create",
			accountLedger{
				"pkey1": testAccount(),
			},
			accountLedger{},
			createErr,
		},
		{
			"return error if there's a db error on save",
			accountLedger{},
			accountLedger{
				"pkey1": testAccount(),
			},
			saveErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockaccountAggStore(ctrl)

			ledger := combineLedgers(tt.new, tt.existing)
			payload := testAccountAggPayload(ledger)

			for key, acnt := range tt.new {
				if tt.result == findByPublicKeyErr {
					dbMock.EXPECT().FindByPublicKey(key).Return(nil, findByPublicKeyErr).Times(1)
					break
				}
				dbMock.EXPECT().FindByPublicKey(key).Return(nil, store.ErrNotFound).Times(1)
				newAccount := newAccountAgg(key, payload.Syncable.Height, payload.Syncable.Time)
				updatedAccount := updateAccountAgg(newAccount, acnt, payload)
				if tt.result == createErr {
					dbMock.EXPECT().Create(updatedAccount).Return(createErr).Times(1)
					break
				}
				dbMock.EXPECT().Create(updatedAccount).Return(nil).Times(1)
			}

			for key, acnt := range tt.existing {
				existingAccount := newAccountAgg(key, 0, *types.NewTimeFromTime(time.Now()))
				dbMock.EXPECT().FindByPublicKey(key).Return(existingAccount, nil).Times(1)
				updatedAcnt := updateAccountAgg(existingAccount, acnt, payload)

				if tt.result == saveErr {
					dbMock.EXPECT().Save(updatedAcnt).Return(saveErr).Times(1)
					break
				}

				dbMock.EXPECT().Save(updatedAcnt).Return(nil).Times(1)
			}

			task := NewAccountAggCreatorTask(dbMock)
			if result := task.Run(ctx, payload); result != tt.result {
				t.Errorf("unexpected result, got: %v; want: %v", nil, result)
				return
			}

			// don't check payload if there was an error
			if tt.result != nil {
				return
			}

			if len(payload.NewAggregatedAccounts) != len(tt.new) {
				t.Errorf("expected payload.NewAggregatedAccounts to contain new accounts, got: %v; want: %v", len(payload.NewAggregatedAccounts), len(tt.new))
				return
			}

			fmt.Printf("payload.UpdatedAggregatedAccounts:%+v\n", payload.UpdatedAggregatedAccounts)
			if len(payload.UpdatedAggregatedAccounts) != len(tt.existing) {
				t.Errorf("expected payload.UpdatedAggregatedAccounts to contain accounts, got: %v; want: %v", len(payload.UpdatedAggregatedAccounts), len(tt.existing))
				return
			}
		})
	}
}

func TestValidatorAggCreatorTask_Run(t *testing.T) {
	setup(t)

	FindByEntityUIDErr := errors.New("FindByEntityUIDErr")
	createErr := errors.New("createErr")
	saveErr := errors.New("saveErr")

	tests := []struct {
		description string
		new         []*validatorpb.Validator
		existing    []*validatorpb.Validator
		parsed      ParsedValidatorsData
		result      error
	}{
		{
			"creates new validators",
			[]*validatorpb.Validator{testValidator("key1")},
			[]*validatorpb.Validator{},
			make(ParsedValidatorsData),
			nil,
		},
		{
			"updates existing accounts",
			[]*validatorpb.Validator{},
			[]*validatorpb.Validator{testValidator("key1")},
			make(ParsedValidatorsData),
			nil,
		},
		{
			"creates and updates accounts",
			[]*validatorpb.Validator{testValidator("key1"), testValidator("key2")},
			[]*validatorpb.Validator{testValidator("key3"), testValidator("key4")},
			make(ParsedValidatorsData),
			nil,
		},
		{
			"return error if there's an unexpected db error on FindByEntityUIDErr",
			[]*validatorpb.Validator{testValidator("key1"), testValidator("key2")},
			[]*validatorpb.Validator{},
			make(ParsedValidatorsData),
			FindByEntityUIDErr,
		},
		{
			"return error if there's a db error on create",
			[]*validatorpb.Validator{testValidator("key1"), testValidator("key2")},
			[]*validatorpb.Validator{},
			make(ParsedValidatorsData),
			createErr,
		},
		{
			"return error if there's a db error on save",
			[]*validatorpb.Validator{},
			[]*validatorpb.Validator{testValidator("key1"), testValidator("key2")},
			make(ParsedValidatorsData),
			saveErr,
		},
		{
			"updates new validators with parsedValidator data",
			[]*validatorpb.Validator{testValidator("key1"), testValidator("key2"), testValidator("key3")},
			[]*validatorpb.Validator{testValidator("key4")},
			ParsedValidatorsData{
				"key1": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: 1,
					TotalShares:          types.NewQuantity(big.NewInt(66)),
				},
				"key2": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 2,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
				"key3": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 0,
					TotalShares:          types.NewQuantity(big.NewInt(68)),
				},
			},
			nil,
		},
		{
			"updates existing validators with parsedValidator data",
			[]*validatorpb.Validator{testValidator("key0")},
			[]*validatorpb.Validator{testValidator("key1"), testValidator("key2"), testValidator("key3")},

			ParsedValidatorsData{
				"key1": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 1,
					TotalShares:          types.NewQuantity(big.NewInt(66)),
				},
				"key2": parsedValidator{
					Proposed:             true,
					PrecommitBlockIdFlag: 2,
					TotalShares:          types.NewQuantity(big.NewInt(67)),
				},
				"key3": parsedValidator{
					Proposed:             false,
					PrecommitBlockIdFlag: 0,
					TotalShares:          types.NewQuantity(big.NewInt(68)),
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockvalidatorAggCreatorStore(ctrl)

			allValidators := append(tt.new, tt.existing...)
			payload := testValidatorAggPayload(allValidators)
			payload.ParsedValidators = tt.parsed

			for _, validator := range tt.new {
				key := validator.GetNode().GetEntityId()
				if tt.result == FindByEntityUIDErr {
					dbMock.EXPECT().FindByEntityUID(key).Return(nil, FindByEntityUIDErr).Times(1)
					break
				}
				dbMock.EXPECT().FindByEntityUID(key).Return(nil, store.ErrNotFound).Times(1)

				newValidator := newValidatorAgg(key, payload.Syncable.Height, payload.Syncable.Time)
				updatedValidator := updateValidatorAgg(newValidator, validator, payload)

				parsedValidator, ok := tt.parsed[key]
				if ok {
					updateParsedValidatorAgg(updatedValidator, parsedValidator, payload)
				}

				if tt.result == createErr {
					dbMock.EXPECT().Create(updatedValidator).Return(createErr).Times(1)
					break
				}
				dbMock.EXPECT().Create(updatedValidator).Return(nil).Times(1)
			}

			for _, raw := range tt.existing {
				key := raw.GetNode().GetEntityId()
				existingValidator := newValidatorAgg(key, 0, *types.NewTimeFromTime(time.Now()))
				dbMock.EXPECT().FindByEntityUID(key).Return(existingValidator, nil).Times(1)
				updated := updateValidatorAgg(existingValidator, raw, payload)

				parsedValidator, ok := tt.parsed[key]
				if ok {
					updateParsedValidatorAgg(updated, parsedValidator, payload)
				}

				if tt.result == saveErr {
					dbMock.EXPECT().Save(updated).Return(saveErr).Times(1)
					break
				}
				dbMock.EXPECT().Save(updated).Return(nil).Times(1)
			}

			task := NewValidatorAggCreatorTask(dbMock)
			if result := task.Run(ctx, payload); result != tt.result {
				t.Errorf("unexpected result, got: %v; want: %v", nil, result)
				return
			}

			// don't check payload if there was an error
			if tt.result != nil {
				return
			}

			if len(payload.NewAggregatedValidators) != len(tt.new) {
				t.Errorf("expected payload.NewAggregatedValidators to contain new validators, got: %v; want: %v", len(payload.NewAggregatedValidators), len(tt.new))
				return
			}

			if len(payload.UpdatedAggregatedValidators) != len(tt.existing) {
				t.Errorf("expected payload.UpdatedAggregatedValidators to contain accounts, got: %v; want: %v", len(payload.UpdatedAggregatedValidators), len(tt.existing))
				return
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

func testValidatorAggPayload(raw []*validatorpb.Validator) *payload {
	return &payload{
		Syncable: &model.Syncable{
			Height: 17,
			Time:   *types.NewTimeFromTime(time.Now()),
		},
		CurrentHeight: 17,
		RawValidators: raw,
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

func testValidator(key string) *validatorpb.Validator {
	return &validatorpb.Validator{
		Address:     randString(5),
		VotingPower: 64,
		Node: &validatorpb.Node{
			EntityId: key,
		},
	}
}

func newValidatorAgg(key string, height int64, _time types.Time) *model.ValidatorAgg {
	return &model.ValidatorAgg{
		Aggregate: &model.Aggregate{
			StartedAtHeight: height,
			StartedAt:       _time,
		},
		EntityUID:     key,
		RecentAddress: randString(5),
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
		RecentAddress:           raw.GetAddress(),
		RecentVotingPower:       raw.GetVotingPower(),
		RecentAsValidatorHeight: pl.Syncable.Height,
	}
}

func updateParsedValidatorAgg(m *model.ValidatorAgg, parsed parsedValidator, pl *payload) {
	m.RecentTotalShares = parsed.TotalShares

	if parsed.PrecommitBlockIdFlag == 1 {
		m.AccumulatedUptimeCount++
	} else if parsed.PrecommitBlockIdFlag == 2 {
		m.AccumulatedUptime++
		m.AccumulatedUptimeCount++
	}

	if parsed.Proposed {
		m.RecentProposedHeight = pl.CurrentHeight
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
