package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/mappers/accountaggmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/entityaggmapper"
	"github.com/figment-networks/oasishub-indexer/models/accountagg"
	"github.com/figment-networks/oasishub-indexer/models/entityagg"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Aggregator interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type aggregator struct {
	accountAggDbRepo accountaggrepo.DbRepo
	entityAggDbRepo  entityaggrepo.DbRepo

	previous *payload
	payloads []*payload
}

func NewAggregator(accountAggDbRepo accountaggrepo.DbRepo, entityAggDbRepo entityaggrepo.DbRepo) *aggregator {
	return &aggregator{
		accountAggDbRepo: accountAggDbRepo,
		entityAggDbRepo:  entityAggDbRepo,
	}
}

func (a *aggregator) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*payload)

	// Aggregate accounts
	ca, ua, err := a.aggregateAccounts(payload)
	if err != nil {
		return nil, err
	}
	payload.NewAggregatedAccounts = ca
	payload.UpdatedAggregatedAccounts = ua

	// Aggregate entities
	ce, ue, err := a.aggregateEntities(payload)
	if err != nil {
		return nil, err
	}
	payload.NewAggregatedEntities = ce
	payload.UpdatedAggregatedEntities = ue

	return payload, nil
}

func (a *aggregator) aggregateAccounts(p *payload) ([]accountagg.Model, []accountagg.Model, errors.ApplicationError) {
	accounts, err := accountaggmapper.ToAggregate(p.StateSyncable)
	if err != nil {
		return nil, nil, err
	}

	var created []accountagg.Model
	var updated []accountagg.Model
	for _, account := range accounts {
		existing, err := a.accountAggDbRepo.GetByPublicKey(account.PublicKey)
		if err != nil {
			if err.Status() == errors.NotFoundError {
				if err := a.accountAggDbRepo.Create(&account); err != nil {
					return nil, nil, err
				}
				created = append(created, account)
			} else {
				return nil, nil, err
			}
		} else {
			existing.UpdateAggAttrs(&account)

			if err := a.accountAggDbRepo.Save(existing); err != nil {
				return nil, nil, err
			}
			updated = append(updated, account)
		}
	}
	return created, updated, nil
}

func (a *aggregator) aggregateEntities(p *payload) ([]entityagg.Model, []entityagg.Model, errors.ApplicationError) {
	entities, err := entityaggmapper.ToAggregate(p.StateSyncable)
	if err != nil {
		return nil, nil, err
	}

	var created []entityagg.Model
	var updated []entityagg.Model
	for _, entity := range entities {
		existing, err := a.entityAggDbRepo.GetByEntityUID(entity.EntityUID)
		if err != nil {
			if err.Status() == errors.NotFoundError {
				if err := a.entityAggDbRepo.Create(&entity); err != nil {
					return nil, nil, err
				}
				created = append(created, entity)
			} else {
				return nil, nil, err
			}
		} else {
			existing.UpdateAggAttrs(entity)

			if err := a.entityAggDbRepo.Save(existing); err != nil {
				return nil, nil, err
			}
			updated = append(updated, entity)
		}
	}
	return created, updated, nil
}