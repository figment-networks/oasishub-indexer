package startpipeline

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/accountagg"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/validatoragg"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatoraggrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

var (
	_ pipeline.AsyncTask = (*AccountAggregateCreator)(nil)
	_ pipeline.AsyncTask = (*ValidatorAggregateCreator)(nil)
)

type AccountAggregateCreator struct {
	accountAggDbRepo accountaggrepo.DbRepo
}

func NewAccountAggregateCreator(accountAggDbRepo accountaggrepo.DbRepo) *AccountAggregateCreator {
	return &AccountAggregateCreator{
		accountAggDbRepo: accountAggDbRepo,
	}
}

func (a *AccountAggregateCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	var created []accountagg.Model
	var updated []accountagg.Model
	for _, accountCalculatedData := range payload.CalculatedAccountsData {
		existing, err := a.accountAggDbRepo.GetByPublicKey(accountCalculatedData.PublicKey)
		if err != nil {
			if err.Status() == errors.NotFoundError {
				accountAgg := &accountagg.Model{
					Aggregate: &shared.Aggregate{
						StartedAtHeight: payload.BlockSyncable.Height,
						StartedAt:       payload.BlockSyncable.Time,
						RecentAtHeight:  payload.CurrentHeight,
						RecentAt:        payload.BlockSyncable.Time,
					},

					PublicKey:                        accountCalculatedData.PublicKey,
					RecentGeneralBalance:             accountCalculatedData.RecentGeneralBalance,
					RecentGeneralNonce:               accountCalculatedData.RecentGeneralNonce,
					RecentEscrowActiveBalance:        accountCalculatedData.RecentEscrowActiveBalance,
					RecentEscrowActiveTotalShares:    accountCalculatedData.RecentEscrowActiveTotalShares,
					RecentEscrowDebondingBalance:     accountCalculatedData.RecentEscrowDebondingBalance,
					RecentEscrowDebondingTotalShares: accountCalculatedData.RecentEscrowDebondingTotalShares,
				}

				if !accountAgg.Valid() {
					errCh <- errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
					return
				}

				if err := a.accountAggDbRepo.Create(accountAgg); err != nil {
					errCh <- err
					return
				}
				created = append(created, *accountAgg)
			} else {
				errCh <- err
				return
			}
		} else {
			accountAgg := &accountagg.Model{
				Aggregate: &shared.Aggregate{
					RecentAtHeight: payload.CurrentHeight,
					RecentAt:       payload.BlockSyncable.Time,
				},

				PublicKey:                        accountCalculatedData.PublicKey,
				RecentGeneralBalance:             accountCalculatedData.RecentGeneralBalance,
				RecentGeneralNonce:               accountCalculatedData.RecentGeneralNonce,
				RecentEscrowActiveBalance:        accountCalculatedData.RecentEscrowActiveBalance,
				RecentEscrowActiveTotalShares:    accountCalculatedData.RecentEscrowActiveTotalShares,
				RecentEscrowDebondingBalance:     accountCalculatedData.RecentEscrowDebondingBalance,
				RecentEscrowDebondingTotalShares: accountCalculatedData.RecentEscrowDebondingTotalShares,
			}

			existing.UpdateAggAttrs(accountAgg)

			if !existing.Valid() {
				errCh <- errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
				return
			}

			if err := a.accountAggDbRepo.Save(existing); err != nil {
				errCh <- err
				return
			}
			updated = append(updated, *accountAgg)
		}
	}
	payload.NewAggregatedAccounts = created
	payload.UpdatedAggregatedAccounts = updated
	errCh <- nil
}

type ValidatorAggregateCreator struct {
	validatorAggDbRepo validatoraggrepo.DbRepo
}

func NewValidatorAggregateCreator(validatorAggDbRepo validatoraggrepo.DbRepo) *ValidatorAggregateCreator {
	return &ValidatorAggregateCreator{
		validatorAggDbRepo: validatorAggDbRepo,
	}
}

func (a *ValidatorAggregateCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	stateRawData, err := syncablemapper.UnmarshalStateData(payload.StateSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}

	var created []validatoragg.Model
	var updated []validatoragg.Model
	for _, entity := range stateRawData.Registry.Entities {
		// check if is validator
		var validatorCalculatedData *CalculatedValidatorData
		for _, d := range payload.CalculatedValidatorsData {
			if d.EntityUID.Equal(types.PublicKey(entity.PublicKey)) {
				validatorCalculatedData = &d
				break
			}
		}

		if validatorCalculatedData != nil {
			existing, err := a.validatorAggDbRepo.GetByEntityUID(types.PublicKey(entity.PublicKey))
			if err != nil {
				if err.Status() == errors.NotFoundError {
					validator := validatoragg.Model{
						Aggregate: &shared.Aggregate{
							StartedAtHeight: payload.CurrentHeight,
							StartedAt:       payload.BlockSyncable.Time,
							RecentAtHeight:  payload.CurrentHeight,
							RecentAt:        payload.BlockSyncable.Time,
						},

						EntityUID:               validatorCalculatedData.EntityUID,
						RecentAddress:           validatorCalculatedData.Address,
						RecentTotalShares:       validatorCalculatedData.TotalShares,
						RecentVotingPower:       validatorCalculatedData.VotingPower,
						RecentAsValidatorHeight: payload.CurrentHeight,
					}

					if validatorCalculatedData.PrecommitValidated == 0 {
						validator.AccumulatedUptime = 0
						validator.AccumulatedUptimeCount = 1
					} else if validatorCalculatedData.PrecommitValidated == 1 {
						validator.AccumulatedUptime = 1
						validator.AccumulatedUptimeCount = 1
					} else {
						// We don't count out of range as offline
						validator.AccumulatedUptime = 0
						validator.AccumulatedUptimeCount = 0
					}

					if validatorCalculatedData.Proposed {
						validator.RecentProposedHeight = payload.CurrentHeight
						validator.AccumulatedProposedCount = 1
					}

					if !validator.Valid() {
						errCh <- errors.NewErrorFromMessage("validator aggregate not valid", errors.NotValid)
						return
					}

					if err := a.validatorAggDbRepo.Create(&validator); err != nil {
						errCh <- err
						return
					}
					created = append(created, validator)
				} else {
					errCh <- err
					return
				}
			} else {
				validator := validatoragg.Model{
					Aggregate: &shared.Aggregate{
						RecentAtHeight: payload.CurrentHeight,
						RecentAt:       payload.StateSyncable.Time,
					},

					RecentAddress:           validatorCalculatedData.Address,
					RecentTotalShares:       validatorCalculatedData.TotalShares,
					RecentVotingPower:       validatorCalculatedData.VotingPower,
					RecentAsValidatorHeight: payload.CurrentHeight,
				}

				if validatorCalculatedData.PrecommitValidated == 0 {
					validator.AccumulatedUptime = existing.AccumulatedUptime
					validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
				} else if validatorCalculatedData.PrecommitValidated == 1 {
					validator.AccumulatedUptime = existing.AccumulatedUptime + 1
					validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
				} else {
					// We don't count out of range as offline
					validator.AccumulatedUptime = existing.AccumulatedUptime
					validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount
				}

				if validatorCalculatedData.Proposed {
					validator.RecentProposedHeight = payload.StateSyncable.Height
					validator.AccumulatedProposedCount = existing.AccumulatedProposedCount + 1
				}

				existing.UpdateAggAttrs(validator)

				if !existing.Valid() {
					errCh <- errors.NewErrorFromMessage("validator aggregate not valid", errors.NotValid)
					return
				}

				if err := a.validatorAggDbRepo.Save(existing); err != nil {
					errCh <- err
					return
				}
				updated = append(updated, validator)
			}
		}
	}
	payload.NewAggregatedValidators = created
	payload.UpdatedAggregatedValidators = updated
	errCh <- nil
}
