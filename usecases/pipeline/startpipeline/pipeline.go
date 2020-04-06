package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/iterators"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Pipeline interface {
	Start(ctx context.Context, iter pipeline.Iterator) Results
}

type processingPipeline struct {
	syncableDbRepo   syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo

	blockSeqDbRepo               blockseqrepo.DbRepo
	validatorSeqDbRepo           validatorseqrepo.DbRepo
	transactionSeqDbRepo         transactionseqrepo.DbRepo
	stakingSeqDbRepo             stakingseqrepo.DbRepo
	accountAggDbRepo             accountaggrepo.DbRepo
	delegationSeqDbRepo          delegationseqrepo.DbRepo
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
	entityAggDbRepo              entityaggrepo.DbRepo

	report   report.Model
	pipeline *pipeline.Pipeline
}

type Details struct {
}

type Results struct {
	SuccessCount int64
	ErrorCount   int64
	Error        *string
	Details      []byte
}

func NewPipeline(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	blockSeqDbRepo blockseqrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	transactionSeqDbRepo transactionseqrepo.DbRepo,
	stakingSeqDbRepo stakingseqrepo.DbRepo,
	accountAggDbRepo accountaggrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo,
	entityAggDbRepo entityaggrepo.DbRepo,
	report report.Model,
) Pipeline {
	// Assemble pipeline
	p := pipeline.New(
		pipeline.FIFO(NewSyncer(syncableDbRepo, syncableProxyRepo)),
		pipeline.FIFO(NewSequencer(blockSeqDbRepo, validatorSeqDbRepo, transactionSeqDbRepo, stakingSeqDbRepo, delegationSeqDbRepo, debondingDelegationSeqDbRepo)),
		pipeline.FIFO(NewAggregator(accountAggDbRepo, entityAggDbRepo)),
	)

	return &processingPipeline{
		syncableDbRepo: syncableDbRepo,

		report:   report,
		pipeline: p,
	}
}

// Process block until the link iterator is exhausted, an error occurs or the
// context is cancelled.
func (c *processingPipeline) Start(ctx context.Context, iter pipeline.Iterator) Results {
	i := iter.(*iterators.HeightIterator)
	source := NewBlockSource(i)
	sink := NewSink(c.syncableDbRepo, c.report)
	err := c.pipeline.Process(ctx, source, sink)
	//TODO: Add details to results
	r := Results{
		SuccessCount: sink.Count(),
		ErrorCount:   i.Length() - sink.Count(),
	}
	if err != nil {
		errMsg := err.Error()
		r.Error = &errMsg
	}
	return r
}
