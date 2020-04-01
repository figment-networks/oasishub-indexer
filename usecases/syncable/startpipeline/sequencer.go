package startpipeline

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/domain/blockdomain"
	"github.com/figment-networks/oasishub-indexer/domain/delegationdomain"
	"github.com/figment-networks/oasishub-indexer/domain/stakingdomain"
	"github.com/figment-networks/oasishub-indexer/domain/transactiondomain"
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/mappers/blockseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/debondingdelegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/delegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/stakingseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/transactionseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/validatorseqmapper"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type Sequencer interface {
	Process(context.Context, pipeline.Payload) (pipeline.Payload, error)
}

type sequencer struct {
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
	payload := p.(*payload)

	// Sequence block
	bs, err := s.sequenceBlock(payload)
	if err != nil {
		return nil, err
	}
	payload.BlockSequence = bs

	// Sequence validators
	vs, err := s.sequenceValidators(payload)
	if err != nil {
		return nil, err
	}
	payload.ValidatorSequences = vs

	// Sequence transaction
	ts, err := s.sequenceTransactions(payload)
	if err != nil {
		return nil, err
	}
	payload.TransactionSequences = ts

	// Sequence staking
	ss, err := s.sequenceStaking(payload)
	if err != nil {
		return nil, err
	}
	payload.StakingSequence = ss

	// Sequence delegations
	sd, err := s.sequenceDelegations(payload)
	if err != nil {
		return nil, err
	}
	payload.DelegationSequences = sd

	// Sequence debonding delegations
	sdd, err := s.sequenceDebondingDelegations(payload)
	if err != nil {
		return nil, err
	}
	payload.DebondingDelegationSequences = sdd

	return payload, nil
}

/*************** Private ***************/

func (s *sequencer) sequenceBlock(p *payload) (*blockdomain.BlockSeq, errors.ApplicationError) {
	sequenced, err := s.blockSeqDbRepo.GetByHeight(p.CurrentHeight)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			toSequence, err := blockseqmapper.FromData(*p.BlockSyncable)
			if err != nil {
				return nil, err
			}
			if err := s.blockSeqDbRepo.Create(toSequence); err != nil {
				return nil, err
			}
			return toSequence, nil
		}
		return nil, err
	}
	return sequenced, nil
}

func (s *sequencer) sequenceValidators(p *payload) ([]*validatordomain.ValidatorSeq, errors.ApplicationError) {
	var res []*validatordomain.ValidatorSeq
	sequenced, err := s.validatorSeqDbRepo.GetByHeight(p.CurrentHeight)
	if err != nil {
		return nil, err
	}

	toSequence, err := validatorseqmapper.FromData(*p.ValidatorsSyncable, *p.BlockSyncable)
	if err != nil {
		return nil, err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		return res, nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		return sequenced, nil
	}

	isSequenced := func(vs *validatordomain.ValidatorSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(*vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.validatorSeqDbRepo.Create(vs); err != nil {
				return nil, err
			}
		}
		res = append(res, vs)

	}

	return res, nil
}

func (s *sequencer) sequenceTransactions(p *payload) ([]*transactiondomain.TransactionSeq, errors.ApplicationError) {
	var res []*transactiondomain.TransactionSeq
	sequenced, err := s.transactionSeqDbRepo.GetByHeight(p.CurrentHeight)
	if err != nil {
		return nil, err
	}

	toSequence, err := transactionseqmapper.FromData(*p.TransactionsSyncable)
	if err != nil {
		return nil, err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		return res, nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		return sequenced, nil
	}

	isSequenced := func(vs *transactiondomain.TransactionSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(*vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.transactionSeqDbRepo.Create(vs); err != nil {
				return nil, err
			}
		}
		res = append(res, vs)

	}

	return res, nil
}

func (s *sequencer) sequenceStaking(p *payload) (*stakingdomain.StakingSeq, errors.ApplicationError) {
	sequenced, err := s.stakingSeqDbRepo.GetByHeight(p.CurrentHeight)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			toSequence, err := stakingseqmapper.FromData(*p.StateSyncable)
			if err != nil {
				return nil, err
			}
			if err := s.stakingSeqDbRepo.Create(toSequence); err != nil {
				return nil, err
			}
			return toSequence, nil
		}
		return nil, err
	}
	return sequenced, nil
}

func (s *sequencer) sequenceDelegations(p *payload) ([]*delegationdomain.DelegationSeq, errors.ApplicationError) {
	var res []*delegationdomain.DelegationSeq
	sequenced, err := s.delegationSeqDbRepo.GetByHeight(p.CurrentHeight)
	if err != nil {
		return nil, err
	}

	toSequence, err := delegationseqmapper.FromData(p.StateSyncable)
	if err != nil {
		return nil, err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		return res, nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		return sequenced, nil
	}

	isSequenced := func(vs *delegationdomain.DelegationSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(*vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.delegationSeqDbRepo.Create(vs); err != nil {
				return nil, err
			}
		}
		res = append(res, vs)
	}
	return res, nil
}

func (s *sequencer) sequenceDebondingDelegations(p *payload) ([]*delegationdomain.DebondingDelegationSeq, errors.ApplicationError) {
	var res []*delegationdomain.DebondingDelegationSeq
	sequenced, err := s.debondingDelegationSeqDbRepo.GetByHeight(p.CurrentHeight)
	if err != nil {
		return nil, err
	}

	toSequence, err := debondingdelegationseqmapper.FromData(p.StateSyncable)
	if err != nil {
		return nil, err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		return res, nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		return sequenced, nil
	}

	isSequenced := func(vs *delegationdomain.DebondingDelegationSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(*vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.debondingDelegationSeqDbRepo.Create(vs); err != nil {
				return nil, err
			}
		}
		res = append(res, vs)
	}
	return res, nil
}
