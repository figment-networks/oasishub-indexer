package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/domain/reportdomain"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Sink interface {
	Consume(context.Context, pipeline.Payload) error
	Count() int64
	Payloads() []*payload
}

type sink struct {
	syncableDbRepo syncablerepo.DbRepo
	report         reportdomain.Report

	count    int64
	payloads []*payload
}

func NewSink(syncableDbRepo syncablerepo.DbRepo, report reportdomain.Report) Sink {
	return &sink{
		syncableDbRepo: syncableDbRepo,
		report:         report,
	}
}

func (s *sink) Consume(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)
	s.count = s.count + 1
	s.payloads = append(s.payloads, payload)

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

func (s *sink) Payloads() []*payload {
	return s.payloads
}
