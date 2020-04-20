package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Sink interface {
	Consume(context.Context, pipeline.Payload) error
	Count() int64
}

type sink struct {
	syncableDbRepo syncablerepo.DbSaver
	report         report.Model

	count    int64
}

func NewSink(syncableDbRepo syncablerepo.DbSaver, report report.Model) Sink {
	return &sink{
		syncableDbRepo: syncableDbRepo,
		report:         report,
	}
}

func (s *sink) Consume(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)
	s.count = s.count + 1

	block := payload.BlockSyncable
	block.MarkProcessed(s.report.ID)
	s.syncableDbRepo.Save(block)

	validators := payload.ValidatorsSyncable
	validators.MarkProcessed(s.report.ID)
	s.syncableDbRepo.Save(validators)

	state := payload.StateSyncable
	state.MarkProcessed(s.report.ID)
	s.syncableDbRepo.Save(state)

	transactions := payload.TransactionsSyncable
	transactions.MarkProcessed(s.report.ID)
	s.syncableDbRepo.Save(transactions)

	return nil
}

func (s *sink) Count() int64 {
	return s.count
}
