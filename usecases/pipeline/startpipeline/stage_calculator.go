package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

// Calculator stage is responsible for calculating and normalizing raw data coming from syncables (syncer stage)
type Calculator interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type calculator struct {
	calculators []pipeline.AsyncTask
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (s *calculator) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	asyncTaskRunner := pipeline.NewAsyncTaskRunner(s.getCalculators())

	return asyncTaskRunner.Run(p)
}

func (s *calculator) getCalculators() []pipeline.AsyncTask{
	if len(s.calculators) == 0 {
		s.calculators = append(s.calculators, NewCalculatedValidatorsData())
		s.calculators = append(s.calculators, NewCalculatedAccountsData())
		s.calculators = append(s.calculators, NewCalculatedEntitiesData())
	}
	return s.calculators
}