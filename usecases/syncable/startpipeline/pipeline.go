package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub/domain/reportdomain"
	"github.com/figment-networks/oasishub/repos/accountaggrepo"
	"github.com/figment-networks/oasishub/repos/blockseqrepo"
	"github.com/figment-networks/oasishub/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub/repos/entityaggrepo"
	"github.com/figment-networks/oasishub/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub/repos/syncablerepo"
	"github.com/figment-networks/oasishub/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub/utils/iterators"
	"github.com/figment-networks/oasishub/utils/pipeline"
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

	report   reportdomain.Report
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
	report reportdomain.Report,
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
