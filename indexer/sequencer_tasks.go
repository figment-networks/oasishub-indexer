package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"reflect"
	"time"
)

var (
	_ pipeline.Task = (*blockSeqCreatorTask)(nil)
	_ pipeline.Task = (*validatorSeqCreatorTask)(nil)
	_ pipeline.Task = (*transactionSeqCreatorTask)(nil)
	_ pipeline.Task = (*stakingSeqCreatorTask)(nil)
	_ pipeline.Task = (*delegationsSeqCreatorTask)(nil)
	_ pipeline.Task = (*debondingDelegationsSeqCreatorTask)(nil)
)

func NewBlockSeqCreatorTask(db *store.Store) *blockSeqCreatorTask {
	return &blockSeqCreatorTask{
		db: db,
	}
}

type blockSeqCreatorTask struct {
	db *store.Store
}

func (t *blockSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("creating block sequence for height %d", payload.CurrentHeight))

	newBlockSeq, err := BlockToSequence(payload.Syncable, payload.RawBlock, payload.ParsedBlock)
	if err != nil {
		return err
	}

	if err := t.db.BlockSeq.CreateIfNotExists(newBlockSeq); err != nil {
		return err
	}

	payload.BlockSequence = newBlockSeq
	return nil
}

func NewValidatorSeqCreatorTask(db *store.Store) *validatorSeqCreatorTask {
	return &validatorSeqCreatorTask{
		db: db,
	}
}

type validatorSeqCreatorTask struct {
	db *store.Store
}

func (t *validatorSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("creating validator sequences for height %d", payload.CurrentHeight))

	var res []model.ValidatorSeq
	sequenced, err := t.db.ValidatorSeq.FindByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	toSequence, err := ValidatorToSequence(payload.Syncable, payload.RawValidators, payload.ParsedValidators)
	if err != nil {
		return err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.ValidatorSequences = res
		return nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.ValidatorSequences = sequenced
		return nil
	}

	isSequenced := func(vs model.ValidatorSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := t.db.ValidatorSeq.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.ValidatorSequences = res
	return nil
}

func NewTransactionSeqCreatorTask(db *store.Store) *transactionSeqCreatorTask {
	return &transactionSeqCreatorTask{
		db: db,
	}
}

type transactionSeqCreatorTask struct {
	db *store.Store
}

func (t *transactionSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("creating transaction sequences for height %d", payload.CurrentHeight))

	var res []model.TransactionSeq
	sequenced, err := t.db.TransactionSeq.FindByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	toSequence, err := TransactionToSequence(payload.Syncable, payload.RawTransactions)
	if err != nil {
		return err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.TransactionSequences = res
		return nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.TransactionSequences = sequenced
		return nil
	}

	isSequenced := func(vs model.TransactionSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := t.db.TransactionSeq.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.TransactionSequences = res
	return nil
}

func NewStakingSeqCreatorTask(db *store.Store) *stakingSeqCreatorTask {
	return &stakingSeqCreatorTask{
		db: db,
	}
}

type stakingSeqCreatorTask struct {
	db *store.Store
}

func (t *stakingSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("creating staking sequence for height %d", payload.CurrentHeight))

	sequenced, err := t.db.StakingSeq.FindByHeight(payload.CurrentHeight)
	if err != nil {
		if err == store.ErrNotFound {
			toSequence, err := StakingToSequence(payload.Syncable, payload.RawState)
			if err != nil {
				return err
			}
			if err := t.db.StakingSeq.Create(toSequence); err != nil {
				return err
			}
			payload.StakingSequence = toSequence
			return nil
		}
		return err
	}
	payload.StakingSequence = sequenced
	return nil
}

type delegationsSeqCreatorTask struct {
	db *store.Store
}

func NewDelegationsSeqCreatorTask(db *store.Store) *delegationsSeqCreatorTask {
	return &delegationsSeqCreatorTask{
		db: db,
	}
}

func (t *delegationsSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("creating delegation sequences for height %d", payload.CurrentHeight))

	var res []model.DelegationSeq
	sequenced, err := t.db.DelegationSeq.FindByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	toSequence, err := DelegationToSequence(payload.Syncable, payload.RawState)
	if err != nil {
		return err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.DelegationSequences = res
		return nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.DelegationSequences = sequenced
		return nil
	}

	isSequenced := func(vs model.DelegationSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := t.db.DelegationSeq.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.DelegationSequences = res
	return nil
}

func NewDebondingDelegationsSeqCreatorTask(db *store.Store) *debondingDelegationsSeqCreatorTask {
	return &debondingDelegationsSeqCreatorTask{
		db: db,
	}
}

type debondingDelegationsSeqCreatorTask struct {
	db *store.Store
}

func (t *debondingDelegationsSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("creating debonding delegation sequences for height %d", payload.CurrentHeight))

	var res []model.DebondingDelegationSeq
	sequenced, err := t.db.DebondingDelegationSeq.FindByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	toSequence, err := DebondingDelegationToSequence(payload.Syncable, payload.RawState)
	if err != nil {
		return err
	}

	// Nothing to sequence
	if len(toSequence) == 0 {
		payload.DebondingDelegationSequences = res
		return nil
	}

	// Everything sequenced and saved to persistence
	if len(sequenced) == len(toSequence) {
		payload.DebondingDelegationSequences = sequenced
		return nil
	}

	isSequenced := func(vs model.DebondingDelegationSeq) bool {
		for _, sv := range sequenced {
			if sv.Equal(vs) {
				return true
			}
		}
		return false
	}

	for _, vs := range toSequence {
		if !isSequenced(vs) {
			if err := t.db.DebondingDelegationSeq.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.DebondingDelegationSequences = res
	return nil
}
