package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Sequencer interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type sequencer struct {
	sequences []pipeline.AsyncTask

	blockSeqDbRepo               blockseqrepo.DbRepo
	validatorSeqDbRepo           validatorseqrepo.DbRepo
	transactionSeqDbRepo         transactionseqrepo.DbRepo
	stakingSeqDbRepo             stakingseqrepo.DbRepo
	delegationSeqDbRepo          delegationseqrepo.DbRepo
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
}

func NewSequencer(
	blockSeqDbRepo blockseqrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	transactionSeqDbRepo transactionseqrepo.DbRepo,
	stakingSeqDbRepo stakingseqrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo,
) Sequencer {
	return &sequencer{
		blockSeqDbRepo:               blockSeqDbRepo,
		validatorSeqDbRepo:           validatorSeqDbRepo,
		transactionSeqDbRepo:         transactionSeqDbRepo,
		stakingSeqDbRepo:             stakingSeqDbRepo,
		delegationSeqDbRepo:          delegationSeqDbRepo,
		debondingDelegationSeqDbRepo: debondingDelegationSeqDbRepo,
	}
}

func (s *sequencer) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	asyncTaskRunner := pipeline.NewAsyncTaskRunner(s.getSequences())

	return asyncTaskRunner.Run(p)
}

func (s *sequencer) getSequences() []pipeline.AsyncTask {
	if len(s.sequences) == 0 {
		s.sequences = append(s.sequences, NewBlockSequenceCreator(s.blockSeqDbRepo))
		s.sequences = append(s.sequences, NewValidatorSequenceCreator(s.validatorSeqDbRepo))
		s.sequences = append(s.sequences, NewTransactionSequencesCreator(s.transactionSeqDbRepo))
		s.sequences = append(s.sequences, NewStakingSequenceCreator(s.stakingSeqDbRepo))
		s.sequences = append(s.sequences, NewDelegationsSequenceCreator(s.delegationSeqDbRepo))
		s.sequences = append(s.sequences, NewDebondingDelegationsSequenceCreator(s.debondingDelegationSeqDbRepo))
	}
	return s.sequences
}
