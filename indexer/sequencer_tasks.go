package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"time"
)

const (
	BlockSeqCreatorTaskName               = "BlockSeqCreator"
	ValidatorSeqCreatorTaskName           = "ValidatorSeqCreator"
	TransactionSeqCreatorTaskName         = "TransactionSeqCreator"
	StakingSeqCreatorTaskName             = "StakingSeqCreator"
	DelegationSeqCreatorTaskName          = "DelegationSeqCreator"
	DebondingDelegationSeqCreatorTaskName = "DebondingDelegationSeqCreator"
)

var (
	_ pipeline.Task = (*blockSeqCreatorTask)(nil)
	_ pipeline.Task = (*validatorSeqCreatorTask)(nil)
	_ pipeline.Task = (*transactionSeqCreatorTask)(nil)
	_ pipeline.Task = (*stakingSeqCreatorTask)(nil)
	_ pipeline.Task = (*delegationSeqCreatorTask)(nil)
	_ pipeline.Task = (*debondingDelegationSeqCreatorTask)(nil)
)

func NewBlockSeqCreatorTask(db BlockSeqStore) *blockSeqCreatorTask {
	return &blockSeqCreatorTask{
		db: db,
	}
}

type blockSeqCreatorTask struct {
	db BlockSeqStore
}

func (t *blockSeqCreatorTask) GetName() string {
	return BlockSeqCreatorTaskName
}

func (t *blockSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	newBlockSeq, err := BlockToSequence(payload.Syncable, payload.RawBlock, payload.ParsedBlock)
	if err != nil {
		return err
	}

	if err := t.db.CreateIfNotExists(newBlockSeq); err != nil {
		return err
	}

	payload.BlockSequence = newBlockSeq
	return nil
}

func NewValidatorSeqCreatorTask(db ValidatorSeqStore) *validatorSeqCreatorTask {
	return &validatorSeqCreatorTask{
		db: db,
	}
}

type validatorSeqCreatorTask struct {
	db ValidatorSeqStore
}

func (t *validatorSeqCreatorTask) GetName() string {
	return ValidatorSeqCreatorTaskName
}

func (t *validatorSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	var res []model.ValidatorSeq
	sequenced, err := t.db.FindByHeight(payload.CurrentHeight)
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
			if err := t.db.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.ValidatorSequences = res
	return nil
}

func NewTransactionSeqCreatorTask(db TransactionSeqStore) *transactionSeqCreatorTask {
	return &transactionSeqCreatorTask{
		db: db,
	}
}

type transactionSeqCreatorTask struct {
	db TransactionSeqStore
}

func (t *transactionSeqCreatorTask) GetName() string {
	return TransactionSeqCreatorTaskName
}

func (t *transactionSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	var res []model.TransactionSeq
	sequenced, err := t.db.FindByHeight(payload.CurrentHeight)
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
			if err := t.db.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.TransactionSequences = res
	return nil
}

func NewStakingSeqCreatorTask(db StakingSeqStore) *stakingSeqCreatorTask {
	return &stakingSeqCreatorTask{
		db: db,
	}
}

type stakingSeqCreatorTask struct {
	db StakingSeqStore
}

func (t *stakingSeqCreatorTask) GetName() string {
	return StakingSeqCreatorTaskName
}

func (t *stakingSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	sequenced, err := t.db.FindByHeight(payload.CurrentHeight)
	if err != nil {
		if err == store.ErrNotFound {
			toSequence, err := StakingToSequence(payload.Syncable, payload.RawState)
			if err != nil {
				return err
			}
			if err := t.db.Create(toSequence); err != nil {
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

type delegationSeqCreatorTask struct {
	db DelegationSeqStore
}

func (t *delegationSeqCreatorTask) GetName() string {
	return DelegationSeqCreatorTaskName
}

func NewDelegationsSeqCreatorTask(db DelegationSeqStore) *delegationSeqCreatorTask {
	return &delegationSeqCreatorTask{
		db: db,
	}
}

func (t *delegationSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	var res []model.DelegationSeq
	sequenced, err := t.db.FindByHeight(payload.CurrentHeight)
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
			if err := t.db.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.DelegationSequences = res
	return nil
}

func NewDebondingDelegationsSeqCreatorTask(db DebondingDelegationSeqStore) *debondingDelegationSeqCreatorTask {
	return &debondingDelegationSeqCreatorTask{
		db: db,
	}
}

type debondingDelegationSeqCreatorTask struct {
	db DebondingDelegationSeqStore
}

func (t *debondingDelegationSeqCreatorTask) GetName() string {
	return DebondingDelegationSeqCreatorTaskName
}

func (t *debondingDelegationSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	var res []model.DebondingDelegationSeq
	sequenced, err := t.db.FindByHeight(payload.CurrentHeight)
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
			if err := t.db.Create(&vs); err != nil {
				return err
			}
		}
		res = append(res, vs)
	}
	payload.DebondingDelegationSequences = res
	return nil
}
