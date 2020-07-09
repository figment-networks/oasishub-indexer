package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/metric"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

const (
	AccountAggCreatorTaskName = "AccountAggCreator"
	ValidatorAggCreatorTaskName = "ValidatorAggCreator"
)

var (
	_ pipeline.Task = (*accountAggCreatorTask)(nil)
	_ pipeline.Task = (*validatorAggCreatorTask)(nil)

	ErrAccountAggNotValid   = errors.New("account aggregator not valid")
	ErrValidatorAggNotValid = errors.New("validator aggregator not valid")
)

func NewAccountAggCreatorTask(db *store.Store) *accountAggCreatorTask {
	return &accountAggCreatorTask{
		db: db,
	}
}

type accountAggCreatorTask struct {
	db *store.Store
}

func (t *accountAggCreatorTask) GetName() string {
	return AccountAggCreatorTaskName
}

func (t *accountAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	var created []model.AccountAgg
	var updated []model.AccountAgg
	for publicKey, rawAccount := range payload.RawState.GetStaking().GetLedger() {
		existing, err := t.db.AccountAgg.FindByPublicKey(publicKey)
		if err != nil {
			if err == store.ErrNotFound {
				accountAgg := &model.AccountAgg{
					Aggregate: &model.Aggregate{
						StartedAtHeight: payload.Syncable.Height,
						StartedAt:       payload.Syncable.Time,
						RecentAtHeight:  payload.Syncable.Height,
						RecentAt:        payload.Syncable.Time,
					},

					PublicKey:                        publicKey,
					RecentGeneralBalance:             types.NewQuantityFromBytes(rawAccount.GetGeneral().GetBalance()),
					RecentGeneralNonce:               rawAccount.GetGeneral().GetNonce(),
					RecentEscrowActiveBalance:        types.NewQuantityFromBytes(rawAccount.GetEscrow().GetActive().GetBalance()),
					RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(rawAccount.GetEscrow().GetActive().GetTotalShares()),
					RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(rawAccount.GetEscrow().GetDebonding().GetBalance()),
					RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(rawAccount.GetEscrow().GetDebonding().GetTotalShares()),
				}

				if !accountAgg.Valid() {
					return ErrAccountAggNotValid
				}

				if err := t.db.AccountAgg.Create(accountAgg); err != nil {
					return err
				}
				created = append(created, *accountAgg)
			} else {
				return err
			}
		} else {
			accountAgg := &model.AccountAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       payload.Syncable.Time,
				},

				PublicKey:                        publicKey,
				RecentGeneralBalance:             types.NewQuantityFromBytes(rawAccount.GetGeneral().GetBalance()),
				RecentGeneralNonce:               rawAccount.GetGeneral().GetNonce(),
				RecentEscrowActiveBalance:        types.NewQuantityFromBytes(rawAccount.GetEscrow().GetActive().GetBalance()),
				RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(rawAccount.GetEscrow().GetActive().GetBalance()),
				RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(rawAccount.GetEscrow().GetDebonding().GetBalance()),
				RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(rawAccount.GetEscrow().GetDebonding().GetTotalShares()),
			}

			existing.UpdateAggAttrs(accountAgg)

			if !existing.Valid() {
				return ErrAccountAggNotValid
			}

			if err := t.db.AccountAgg.Save(existing); err != nil {
				return err
			}
			updated = append(updated, *accountAgg)
		}
	}
	payload.NewAggregatedAccounts = created
	payload.UpdatedAggregatedAccounts = updated
	return nil
}

func NewValidatorAggCreatorTask(db *store.Store) *validatorAggCreatorTask {
	return &validatorAggCreatorTask{
		db: db,
	}
}

type validatorAggCreatorTask struct {
	db *store.Store
}

func (t *validatorAggCreatorTask) GetName() string {
	return ValidatorAggCreatorTaskName
}

func (t *validatorAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	var created []model.ValidatorAgg
	var updated []model.ValidatorAgg
	for _, rawValidator := range payload.RawValidators {
		existing, err := t.db.ValidatorAgg.FindByEntityUID(rawValidator.GetNode().GetEntityId())
		if err != nil {
			if err == store.ErrNotFound {
				// Create new
				validator := model.ValidatorAgg{
					Aggregate: &model.Aggregate{
						StartedAtHeight: payload.Syncable.Height,
						StartedAt:       payload.Syncable.Time,
						RecentAtHeight:  payload.Syncable.Height,
						RecentAt:        payload.Syncable.Time,
					},

					EntityUID:               rawValidator.GetNode().EntityId,
					RecentAddress:           rawValidator.GetAddress(),
					RecentVotingPower:       rawValidator.GetVotingPower(),
					RecentAsValidatorHeight: payload.Syncable.Height,
				}

				parsedValidator, ok := payload.ParsedValidators[rawValidator.GetNode().GetEntityId()]
				if ok {
					validator.RecentTotalShares = parsedValidator.TotalShares

					if parsedValidator.PrecommitBlockIdFlag == 1 {
						// Not validated
						validator.AccumulatedUptime = 0
						validator.AccumulatedUptimeCount = 1
					} else if parsedValidator.PrecommitBlockIdFlag == 2 {
						// Validated
						validator.AccumulatedUptime = 1
						validator.AccumulatedUptimeCount = 1
					} else {
						// Nil validated
						validator.AccumulatedUptime = 0
						validator.AccumulatedUptimeCount = 0
					}

					if parsedValidator.Proposed {
						validator.RecentProposedHeight = payload.CurrentHeight
						validator.AccumulatedProposedCount = 1
					}
				}

				if !validator.Valid() {
					return ErrValidatorAggNotValid
				}

				if err := t.db.ValidatorAgg.Create(&validator); err != nil {
					return err
				}
				created = append(created, validator)
			} else {
				return err
			}
		} else {
			// Update
			validator := model.ValidatorAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       payload.Syncable.Time,
				},

				RecentAddress:           rawValidator.GetAddress(),
				RecentVotingPower:       rawValidator.GetVotingPower(),
				RecentAsValidatorHeight: payload.Syncable.Height,
			}

			parsedValidator, ok := payload.ParsedValidators[rawValidator.GetNode().GetEntityId()]
			if ok {
				validator.RecentTotalShares = parsedValidator.TotalShares

				if parsedValidator.PrecommitBlockIdFlag == 1 {
					// Not validated
					validator.AccumulatedUptime = existing.AccumulatedUptime
					validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
				} else if parsedValidator.PrecommitBlockIdFlag == 2 {
					// Validated
					validator.AccumulatedUptime = existing.AccumulatedUptime + 1
					validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
				} else {
					// Validated nil
					validator.AccumulatedUptime = existing.AccumulatedUptime
					validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount
				}

				if parsedValidator.Proposed {
					validator.RecentProposedHeight = payload.Syncable.Height
					validator.AccumulatedProposedCount = existing.AccumulatedProposedCount + 1
				}
			}

			existing.UpdateAggAttrs(validator)

			if !existing.Valid() {
				return ErrValidatorAggNotValid
			}

			if err := t.db.ValidatorAgg.Save(existing); err != nil {
				return err
			}
			updated = append(updated, validator)
		}
	}
	payload.NewAggregatedValidators = created
	payload.UpdatedAggregatedValidators = updated
	return nil
}
