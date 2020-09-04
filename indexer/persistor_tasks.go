package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	TaskNameBalanceEventPersistor = "BalanceEventPersistor"
	TaskNameSyncerPersistor       = "SyncerPersistor"
	TaskNameBlockSeqPersistor     = "BlockSeqPersistor"
	TaskNameValidatorSeqPersistor = "ValidatorSeqPersistor"
	TaskNameValidatorAggPersistor = "ValidatorAggPersistor"
	TaskNameSystemEventPersistor  = "SystemEventPersistor"
)

func NewSyncerPersistorTask(db SyncerPersistorTaskStore) pipeline.Task {
	return &syncerPersistorTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameSyncerPersistor),
	}
}

type SyncerPersistorTaskStore interface {
	CreateOrUpdate(val *model.Syncable) error
}

type syncerPersistorTask struct {
	db             SyncerPersistorTaskStore
	metricObserver metrics.Observer
}

func (t *syncerPersistorTask) GetName() string {
	return TaskNameSyncerPersistor
}

func (t *syncerPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return t.db.CreateOrUpdate(payload.Syncable)
}

func NewBlockSeqPersistorTask(db BlockSeqPersistorTaskStore) pipeline.Task {
	return &blockSeqPersistorTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameBlockSeqPersistor),
	}
}

type blockSeqPersistorTask struct {
	db             BlockSeqPersistorTaskStore
	metricObserver metrics.Observer
}

type BlockSeqPersistorTaskStore interface {
	Create(record interface{}) error
	Save(record interface{}) error
}

func (t *blockSeqPersistorTask) GetName() string {
	return TaskNameBlockSeqPersistor
}

func (t *blockSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	if payload.NewBlockSequence != nil {
		return t.db.Create(payload.NewBlockSequence)
	}

	if payload.UpdatedBlockSequence != nil {
		return t.db.Save(payload.UpdatedBlockSequence)
	}

	return nil
}

func NewValidatorSeqPersistorTask(db ValidatorSeqPersistorTaskStore) pipeline.Task {
	return &validatorSeqPersistorTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameValidatorSeqPersistor),
	}
}

type ValidatorSeqPersistorTaskStore interface {
	Create(record interface{}) error
	Save(record interface{}) error
}

type validatorSeqPersistorTask struct {
	db             ValidatorSeqPersistorTaskStore
	metricObserver metrics.Observer
}

func (t *validatorSeqPersistorTask) GetName() string {
	return TaskNameValidatorSeqPersistor
}

func (t *validatorSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, sequence := range payload.NewValidatorSequences {
		if err := t.db.Create(&sequence); err != nil {
			return err
		}
	}

	for _, sequence := range payload.UpdatedValidatorSequences {
		if err := t.db.Save(&sequence); err != nil {
			return err
		}
	}

	return nil
}

func NewValidatorAggPersistorTask(db ValidatorAggPersistorTaskStore) pipeline.Task {
	return &validatorAggPersistorTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameValidatorAggPersistor),
	}
}

type ValidatorAggPersistorTaskStore interface {
	Create(record interface{}) error
	Save(record interface{}) error
}

type validatorAggPersistorTask struct {
	db             ValidatorAggPersistorTaskStore
	metricObserver metrics.Observer
}

func (t *validatorAggPersistorTask) GetName() string {
	return TaskNameValidatorAggPersistor
}

func (t *validatorAggPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, aggregate := range payload.NewAggregatedValidators {
		if err := t.db.Create(&aggregate); err != nil {
			return err
		}
	}

	for _, aggregate := range payload.UpdatedAggregatedValidators {
		if err := t.db.Save(&aggregate); err != nil {
			return err
		}
	}

	return nil
}

func NewSystemEventPersistorTask(db SystemEventPersistorTaskStore) pipeline.Task {
	return &systemEventPersistorTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameSystemEventPersistor),
	}
}

type SystemEventPersistorTaskStore interface {
	CreateOrUpdate(*model.SystemEvent) error
}

type systemEventPersistorTask struct {
	db             SystemEventPersistorTaskStore
	metricObserver metrics.Observer
}

func (t *systemEventPersistorTask) GetName() string {
	return TaskNameSystemEventPersistor
}

func (t *systemEventPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, systemEvent := range payload.SystemEvents {
		if err := t.db.CreateOrUpdate(systemEvent); err != nil {
			return err
		}
	}

	return nil
}

func NewBalanceEventPersistorTask(db BalanceEventPersistorTaskStore) pipeline.Task {
	return &balanceEventPersistorTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameBalanceEventPersistor),
	}
}

type BalanceEventPersistorTaskStore interface {
	CreateOrUpdate(*model.BalanceEvent) error
}

type balanceEventPersistorTask struct {
	db             BalanceEventPersistorTaskStore
	metricObserver metrics.Observer
}

func (t *balanceEventPersistorTask) GetName() string {
	return TaskNameBalanceEventPersistor
}

func (t *balanceEventPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, balanceEvent := range payload.BalanceEvents {
		if err := t.db.CreateOrUpdate(&balanceEvent); err != nil {
			return err
		}
	}

	return nil
}
