package startpipeline

import (
	"github.com/figment-networks/oasishub-indexer/mappers/blockseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/debondingdelegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/delegationseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/stakingseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/transactionseqmapper"
	"github.com/figment-networks/oasishub-indexer/mappers/validatorseqmapper"
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/models/transactionseq"
	"github.com/figment-networks/oasishub-indexer/models/validatorseq"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

var (
	_ pipeline.AsyncTask = (*BlockSequenceCreator)(nil)
	_ pipeline.AsyncTask = (*ValidatorSequenceCreator)(nil)
	_ pipeline.AsyncTask = (*TransactionSequencesCreator)(nil)
	_ pipeline.AsyncTask = (*StakingSequenceCreator)(nil)
	_ pipeline.AsyncTask = (*DelegationsSequenceCreator)(nil)
	_ pipeline.AsyncTask = (*DebondingDelegationsSequenceCreator)(nil)
)

type BlockSequenceCreator struct {
	blockSeqDbRepo blockseqrepo.DbRepo
}

func NewBlockSequenceCreator(blockSeqDbRepo blockseqrepo.DbRepo) *BlockSequenceCreator {
	return &BlockSequenceCreator{
		blockSeqDbRepo: blockSeqDbRepo,
	}
}

func (s *BlockSequenceCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	sequenced, err := s.blockSeqDbRepo.GetByHeight(payload.CurrentHeight)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			toSequence, err := blockseqmapper.ToSequence(*payload.BlockSyncable, *payload.ValidatorsSyncable)
			if err != nil {
				errCh <- err
				return
			}
			if err := s.blockSeqDbRepo.Create(toSequence); err != nil {
				errCh <- err
				return
			}
			payload.BlockSequence = sequenced
			errCh <- nil
			return
		}
		errCh <- err
		return
	}
	payload.BlockSequence = sequenced
	errCh <- nil
}

type ValidatorSequenceCreator struct {
	validatorSeqDbRepo validatorseqrepo.DbRepo
}

func NewValidatorSequenceCreator(validatorSeqDbRepo validatorseqrepo.DbRepo) *ValidatorSequenceCreator {
	return &ValidatorSequenceCreator{
		validatorSeqDbRepo: validatorSeqDbRepo,
	}
}

func (s *ValidatorSequenceCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	var res []validatorseq.Model
	sequenced, err := s.validatorSeqDbRepo.GetByHeight(payload.CurrentHeight)
	if err != nil {
		errCh <- err
		return
	}

	toSequence, err := validatorseqmapper.ToSequence(*payload.ValidatorsSyncable, *payload.BlockSyncable, *payload.StateSyncable)
	if err != nil {
		errCh <- err
		return
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.ValidatorSequences = res
		errCh <- nil
		return
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.ValidatorSequences = sequenced
		errCh <- nil
		return
	}

	isSequenced := func(vs validatorseq.Model) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.validatorSeqDbRepo.Create(&vs); err != nil {
				errCh <- err
			}
		}
		res = append(res, vs)
	}
	payload.ValidatorSequences = res
	errCh <- nil
}

type TransactionSequencesCreator struct {
	transactionSeqDbRepo transactionseqrepo.DbRepo
}

func NewTransactionSequencesCreator(transactionSeqDbRepo transactionseqrepo.DbRepo) *TransactionSequencesCreator {
	return &TransactionSequencesCreator{
		transactionSeqDbRepo: transactionSeqDbRepo,
	}
}

func (s *TransactionSequencesCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	var res []transactionseq.Model
	sequenced, err := s.transactionSeqDbRepo.GetByHeight(payload.CurrentHeight)
	if err != nil {
		errCh <- err
		return
	}

	toSequence, err := transactionseqmapper.ToSequence(*payload.TransactionsSyncable)
	if err != nil {
		errCh <- err
		return
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.TransactionSequences = res
		errCh <- nil
		return
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.TransactionSequences = sequenced
		errCh <- nil
		return
	}

	isSequenced := func(vs transactionseq.Model) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.transactionSeqDbRepo.Create(&vs); err != nil {
				errCh <- err
			}
		}
		res = append(res, vs)
	}
	payload.TransactionSequences = res
	errCh <- nil
}

type StakingSequenceCreator struct {
	stakingSeqDbRepo stakingseqrepo.DbRepo
}

func NewStakingSequenceCreator(stakingSeqDbRepo stakingseqrepo.DbRepo) *StakingSequenceCreator {
	return &StakingSequenceCreator{
		stakingSeqDbRepo: stakingSeqDbRepo,
	}
}

func (s *StakingSequenceCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	sequenced, err := s.stakingSeqDbRepo.GetByHeight(payload.CurrentHeight)
	if err != nil {
		if err.Status() == errors.NotFoundError {
			toSequence, err := stakingseqmapper.ToSequence(*payload.StateSyncable)
			if err != nil {
				errCh <- err
				return
			}
			if err := s.stakingSeqDbRepo.Create(toSequence); err != nil {
				errCh <- err
				return
			}
			payload.StakingSequence = toSequence
			errCh <- nil
			return
		}
		errCh <- err
		return
	}
	payload.StakingSequence = sequenced
	errCh <- nil
}

type DelegationsSequenceCreator struct {
	delegationSeqDbRepo delegationseqrepo.DbRepo
}

func NewDelegationsSequenceCreator(delegationSeqDbRepo delegationseqrepo.DbRepo) *DelegationsSequenceCreator {
	return &DelegationsSequenceCreator{
		delegationSeqDbRepo: delegationSeqDbRepo,
	}
}

func (s *DelegationsSequenceCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	var res []delegationseq.Model
	sequenced, err := s.delegationSeqDbRepo.GetByHeight(payload.CurrentHeight)
	if err != nil {
		errCh <- err
		return
	}

	toSequence, err := delegationseqmapper.ToSequence(payload.StateSyncable)
	if err != nil {
		errCh <- err
		return
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.DelegationSequences = res
		errCh <- nil
		return
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.DelegationSequences = sequenced
		errCh <- nil
		return
	}

	isSequenced := func(vs delegationseq.Model) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.delegationSeqDbRepo.Create(&vs); err != nil {
				errCh <- err
				return
			}
		}
		res = append(res, vs)
	}
	payload.DelegationSequences = res
	errCh <- nil
}

type DebondingDelegationsSequenceCreator struct {
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
}

func NewDebondingDelegationsSequenceCreator(debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo) *DebondingDelegationsSequenceCreator {
	return &DebondingDelegationsSequenceCreator{
		debondingDelegationSeqDbRepo: debondingDelegationSeqDbRepo,
	}
}

func (s *DebondingDelegationsSequenceCreator) Run(errCh chan<- error, p pipeline.Payload) {
	payload := p.(*payload)
	var res []debondingdelegationseq.Model
	sequenced, err := s.debondingDelegationSeqDbRepo.GetByHeight(payload.CurrentHeight)
	if err != nil {
		errCh <- err
		return
	}

	toSequence, err := debondingdelegationseqmapper.ToSequence(payload.StateSyncable)
	if err != nil {
		errCh <- err
		return
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.DebondingDelegationSequences = res
		errCh <- nil
		return
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.DebondingDelegationSequences = sequenced
		errCh <- nil
		return
	}

	isSequenced := func(vs debondingdelegationseq.Model) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := s.debondingDelegationSeqDbRepo.Create(&vs); err != nil {
				errCh <- err
			}
		}
		res = append(res, vs)
	}
	payload.DebondingDelegationSequences = res
	errCh <- nil
}
