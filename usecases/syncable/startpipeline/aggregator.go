package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/domain/accountdomain"
	"github.com/figment-networks/oasishub-indexer/domain/entitydomain"
	"github.com/figment-networks/oasishub-indexer/mappers/accountaggmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/entityaggmapper"
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

func (a *aggregator) aggregateAccounts(p *payload) ([]*accountdomain.AccountAgg, []*accountdomain.AccountAgg, errors.ApplicationError) {
	accounts, err := accountaggmapper.FromData(p.StateSyncable)
	if err != nil {
		return nil, nil, err
	}

	var created []*accountdomain.AccountAgg
	var updated []*accountdomain.AccountAgg
	for _, account := range accounts {
		existing, err := a.accountAggDbRepo.GetByPublicKey(account.PublicKey)
		if err != nil {
			if err.Status() == errors.NotFoundError {
				if err := a.accountAggDbRepo.Create(account); err != nil {
					return nil, nil, err
				}
				created = append(created, account)
			} else {
				return nil, nil, err
			}
		} else {
			existing.Update(account)

			if err := a.accountAggDbRepo.Save(existing); err != nil {
				return nil, nil, err
			}
			updated = append(updated, account)
		}
	}
	return created, updated, nil
}

func (a *aggregator) aggregateEntities(p *payload) ([]*entitydomain.EntityAgg, []*entitydomain.EntityAgg, errors.ApplicationError) {
	entities, err := entityaggmapper.FromData(p.StateSyncable)
	if err != nil {
		return nil, nil, err
	}

	var created []*entitydomain.EntityAgg
	var updated []*entitydomain.EntityAgg
	for _, entity := range entities {
		existing, err := a.entityAggDbRepo.GetByEntityUID(entity.EntityUID)
		if err != nil {
			if err.Status() == errors.NotFoundError {
				if err := a.entityAggDbRepo.Create(entity); err != nil {
					return nil, nil, err
				}
				created = append(created, entity)
			} else {
				return nil, nil, err
			}
		} else {
			existing.Update(entity)

			if err := a.entityAggDbRepo.Save(existing); err != nil {
				return nil, nil, err
			}
			updated = append(updated, entity)
		}
	}
	return created, updated, nil
}