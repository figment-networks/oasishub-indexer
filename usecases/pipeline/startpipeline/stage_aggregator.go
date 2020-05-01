package startpipeline

import (
	"context"
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

type Aggregator interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type aggregator struct {
	accountAggDbRepo   accountaggrepo.DbRepo
	validatorAggDbRepo validatoraggrepo.DbRepo

	previous *payload
	payloads []*payload
}

func NewAggregator(accountAggDbRepo accountaggrepo.DbRepo, validatorAggDbRepo validatoraggrepo.DbRepo) *aggregator {
	return &aggregator{
		accountAggDbRepo:   accountAggDbRepo,
		validatorAggDbRepo: validatorAggDbRepo,
	}
}

func (a *aggregator) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*payload)

	if err := a.aggregateAccounts(payload); err != nil {
		return nil, err
	}

	if err := a.aggregateValidators(payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *aggregator) aggregateAccounts(p *payload) errors.ApplicationError {
	accountsCalculatedData, err := CalculateAccountsData(p.StateSyncable)
	if err != nil {
		return err
	}

	var created []accountagg.Model
	var updated []accountagg.Model
	for _, accountCalculatedData := range accountsCalculatedData {
		existing, err := a.accountAggDbRepo.GetByPublicKey(accountCalculatedData.PublicKey)
		if err != nil {
			if err.Status() == errors.NotFoundError {
				accountAgg := &accountagg.Model{
					Aggregate: &shared.Aggregate{
						StartedAtHeight: p.BlockSyncable.Height,
						StartedAt:       p.BlockSyncable.Time,
						RecentAtHeight:  p.CurrentHeight,
						RecentAt:        p.BlockSyncable.Time,
					},

					PublicKey:                         accountCalculatedData.PublicKey,
					RecentGeneralBalance:             accountCalculatedData.RecentGeneralBalance,
					RecentGeneralNonce:               accountCalculatedData.RecentGeneralNonce,
					RecentEscrowActiveBalance:        accountCalculatedData.RecentEscrowActiveBalance,
					RecentEscrowActiveTotalShares:    accountCalculatedData.RecentEscrowActiveTotalShares,
					RecentEscrowDebondingBalance:     accountCalculatedData.RecentEscrowDebondingBalance,
					RecentEscrowDebondingTotalShares: accountCalculatedData.RecentEscrowDebondingTotalShares,
				}

				if !accountAgg.Valid() {
					return errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
				}

				if err := a.accountAggDbRepo.Create(accountAgg); err != nil {
					return err
				}
				created = append(created, *accountAgg)
			} else {
				return err
			}
		} else {
			accountAgg := &accountagg.Model{
				Aggregate: &shared.Aggregate{
					RecentAtHeight:  p.CurrentHeight,
					RecentAt:        p.BlockSyncable.Time,
				},

				PublicKey:                         accountCalculatedData.PublicKey,
				RecentGeneralBalance:             accountCalculatedData.RecentGeneralBalance,
				RecentGeneralNonce:               accountCalculatedData.RecentGeneralNonce,
				RecentEscrowActiveBalance:        accountCalculatedData.RecentEscrowActiveBalance,
				RecentEscrowActiveTotalShares:    accountCalculatedData.RecentEscrowActiveTotalShares,
				RecentEscrowDebondingBalance:     accountCalculatedData.RecentEscrowDebondingBalance,
				RecentEscrowDebondingTotalShares: accountCalculatedData.RecentEscrowDebondingTotalShares,
			}

			existing.UpdateAggAttrs(accountAgg)

			if !existing.Valid() {
				return errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
			}

			if err := a.accountAggDbRepo.Save(existing); err != nil {
				return err
			}
			updated = append(updated, *accountAgg)
		}
	}
	p.NewAggregatedAccounts = created
	p.UpdatedAggregatedAccounts = updated
	return nil
}

func (a *aggregator) aggregateValidators(p *payload) errors.ApplicationError {
	stateRawData, err := syncablemapper.UnmarshalStateData(p.StateSyncable.Data)
	if err != nil {
		return err
	}
	validatorsCalculatedData, err := CalculateValidatorsData(p.ValidatorsSyncable, p.BlockSyncable, p.StateSyncable)
	if err != nil {
		return err
	}

	var created []validatoragg.Model
	var updated []validatoragg.Model
	for _, entity := range stateRawData.Registry.Entities {
		// check if is validator
		var validatorCalculatedData *CalculatedValidatorData
		for _, d := range validatorsCalculatedData {
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
							StartedAtHeight: p.CurrentHeight,
							StartedAt:       p.BlockSyncable.Time,
							RecentAtHeight:  p.CurrentHeight,
							RecentAt:        p.BlockSyncable.Time,
						},

						EntityUID:               validatorCalculatedData.EntityUID,
						RecentAddress:           validatorCalculatedData.Address,
						RecentTotalShares:       validatorCalculatedData.TotalShares,
						RecentVotingPower:       validatorCalculatedData.VotingPower,
						RecentAsValidatorHeight: p.CurrentHeight,
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
						validator.RecentProposedHeight = p.CurrentHeight
						validator.AccumulatedProposedCount = 1
					}

					if !validator.Valid() {
						return errors.NewErrorFromMessage("validator aggregate not valid", errors.NotValid)
					}

					if err := a.validatorAggDbRepo.Create(&validator); err != nil {
						return err
					}
					created = append(created, validator)
				} else {
					return err
				}
			} else {
				validator := validatoragg.Model{
					Aggregate: &shared.Aggregate{
						RecentAtHeight:  p.CurrentHeight,
						RecentAt:        p.StateSyncable.Time,
					},

					RecentAddress:           validatorCalculatedData.Address,
					RecentTotalShares:       validatorCalculatedData.TotalShares,
					RecentVotingPower:       validatorCalculatedData.VotingPower,
					RecentAsValidatorHeight: p.CurrentHeight,
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
					validator.RecentProposedHeight = p.StateSyncable.Height
					validator.AccumulatedProposedCount = existing.AccumulatedProposedCount + 1
				}

				existing.UpdateAggAttrs(validator)

				if !existing.Valid() {
					return errors.NewErrorFromMessage("validator aggregate not valid", errors.NotValid)
				}

				if err := a.validatorAggDbRepo.Save(existing); err != nil {
					return err
				}
				updated = append(updated, validator)
			}
		}
	}
	p.NewAggregatedValidators = created
	p.UpdatedAggregatedValidators = updated
	return nil
}
