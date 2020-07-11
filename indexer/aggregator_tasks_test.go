package indexer

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
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

func testAccountAggPayload(ledger accountLedger) *payload {
	return &payload{
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
	tokens := make([][]byte, 5)
	for i := range tokens {
		token := make([]byte, 10)
		rand.Read(token)
		tokens[i] = token
	}

	return &accountpb.Account{
		General: &accountpb.GeneralAccount{
			Balance: tokens[0],
			Nonce:   100,
		},
		Escrow: &accountpb.EscrowAccount{
			Active: &accountpb.SharePool{
				Balance:     tokens[1],
				TotalShares: tokens[2],
			},
			Debonding: &accountpb.SharePool{
				Balance:     tokens[3],
				TotalShares: tokens[4],
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

func updateAccountAgg(acntag *model.AccountAgg, acnt *accountpb.Account, pl *payload) *model.AccountAgg {
	m := &model.AccountAgg{
		Aggregate: &model.Aggregate{
			StartedAtHeight: acntag.Aggregate.StartedAtHeight,
			StartedAt:       acntag.Aggregate.StartedAt,
			RecentAtHeight:  pl.Syncable.Height,
			RecentAt:        pl.Syncable.Time,
		},
		PublicKey: acntag.PublicKey,

		RecentGeneralBalance:             types.NewQuantityFromBytes(acnt.GetGeneral().GetBalance()),
		RecentGeneralNonce:               acnt.GetGeneral().GetNonce(),
		RecentEscrowActiveBalance:        types.NewQuantityFromBytes(acnt.GetEscrow().GetActive().GetBalance()),
		RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(acnt.GetEscrow().GetActive().GetTotalShares()),
		RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(acnt.GetEscrow().GetDebonding().GetBalance()),
		RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(acnt.GetEscrow().GetDebonding().GetTotalShares()),
	}
	return m
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
