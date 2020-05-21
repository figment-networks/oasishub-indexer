package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatoraggrepo"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Aggregator interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type aggregator struct {
	aggregates []pipeline.AsyncTask

	accountAggDbRepo   accountaggrepo.DbRepo
	validatorAggDbRepo validatoraggrepo.DbRepo
}

func NewAggregator(accountAggDbRepo accountaggrepo.DbRepo, validatorAggDbRepo validatoraggrepo.DbRepo) *aggregator {
	return &aggregator{
		accountAggDbRepo:   accountAggDbRepo,
		validatorAggDbRepo: validatorAggDbRepo,
	}
}

func (a *aggregator) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	asyncTaskRunner := pipeline.NewAsyncTaskRunner(a.getAggregates())

	return asyncTaskRunner.Run(p)
}

func (a *aggregator) getAggregates() []pipeline.AsyncTask {
	if len(a.aggregates) == 0 {
		a.aggregates = append(a.aggregates, NewAccountAggregateCreator(a.accountAggDbRepo))
		a.aggregates = append(a.aggregates, NewValidatorAggregateCreator(a.validatorAggDbRepo))
	}
	return a.aggregates
}
