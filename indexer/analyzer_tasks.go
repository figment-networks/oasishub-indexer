package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
	"math"
	"time"
)

const (
	SystemEventCreatorTaskName = "SystemEventCreator"
)

var (
	ErrActiveEscrowBalanceOutsideOfRange = errors.New("active escrow balance is outside of specified buckets")
	ErrCommissionOutsideOfRange          = errors.New("commission is outside of specified buckets")

	MaxValidatorSequences int64 = 1000
	MissedForMaxThreshold int64 = 50
	MissedInRowThreshold  int64 = 50
)

// NewSystemEventCreatorTask creates system events
func NewSystemEventCreatorTask(cfg *config.Config, s SystemEventCreatorStore) *systemEventCreatorTask {
	return &systemEventCreatorTask{
		cfg: cfg,

		SystemEventCreatorStore: s,
	}
}

type SystemEventCreatorStore interface {
	FindByHeight(int64) ([]model.ValidatorSeq, error)
	FindLastByAddress(string, int64) ([]model.ValidatorSeq, error)
}

type systemEventCreatorTask struct {
	cfg *config.Config

	SystemEventCreatorStore
}

type systemEventRawData map[string]interface{}

func (t *systemEventCreatorTask) GetName() string {
	return SystemEventCreatorTaskName
}

func (t *systemEventCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", "Analyzer", t.GetName(), payload.CurrentHeight))

	currHeightValidatorSequences := append(payload.NewValidatorSequences, payload.UpdatedValidatorSequences...)
	prevHeightValidatorSequences, err := t.getPrevHeightValidatorSequences(payload)
	if err != nil {
		return err
	}

	valueChangeSystemEvents, err := t.getValueChangeSystemEvents(currHeightValidatorSequences, prevHeightValidatorSequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, valueChangeSystemEvents...)

	activeSetPresenceChangeSystemEvents, err := t.getActiveSetPresenceChangeSystemEvents(currHeightValidatorSequences, prevHeightValidatorSequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, activeSetPresenceChangeSystemEvents...)

	missedBlocksSystemEvents, err := t.getMissedBlocksSystemEvents(currHeightValidatorSequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, missedBlocksSystemEvents...)

	return nil
}

func (t *systemEventCreatorTask) getPrevHeightValidatorSequences(payload *payload) ([]model.ValidatorSeq, error) {
	var prevHeightValidatorSequences []model.ValidatorSeq
	if payload.CurrentHeight > t.cfg.FirstBlockHeight {
		var err error
		prevHeightValidatorSequences, err = t.SystemEventCreatorStore.FindByHeight(payload.CurrentHeight - 1)
		if err != nil {
			if err != store.ErrNotFound {
				return nil, err
			}
		}
	}
	return prevHeightValidatorSequences, nil
}

func (t *systemEventCreatorTask) getMissedBlocksSystemEvents(currHeightValidatorSequences []model.ValidatorSeq) ([]*model.SystemEvent, error) {
	var systemEvents []*model.SystemEvent
	for _, validatorSequence := range currHeightValidatorSequences {
		// When current height validator has validated the block no need to check last records
		if t.isValidated(validatorSequence) {
			return systemEvents, nil
		}

		lastValidatorSequencesForAddress, err := t.SystemEventCreatorStore.FindLastByAddress(validatorSequence.Address, MaxValidatorSequences)
		if err != nil {
			if err == store.ErrNotFound {
				return systemEvents, nil
			} else {
				return nil, err
			}
		} else {
			var validatorSequencesToCheck []model.ValidatorSeq
			validatorSequencesToCheck = append([]model.ValidatorSeq{validatorSequence}, lastValidatorSequencesForAddress...)
			totalMissedCount := t.getTotalMissed(validatorSequencesToCheck)

			logger.Debug(fmt.Sprintf("total missed blocks for last %d blocks for address %s: %d", MaxValidatorSequences, validatorSequence.Address, totalMissedCount))

			if totalMissedCount == MissedForMaxThreshold {
				newSystemEvent, err := t.newSystemEvent(validatorSequence, model.SystemEventMissedMofN, systemEventRawData{
					"threshold":               MissedForMaxThreshold,
					"max_validator_sequences": MaxValidatorSequences,
				})
				if err != nil {
					return nil, err
				}

				systemEvents = append(systemEvents, newSystemEvent)
			}

			missedInRowCount := t.getMissedInRow(validatorSequencesToCheck, MissedInRowThreshold)

			logger.Debug(fmt.Sprintf("total missed blocks in a row for address %s: %d", validatorSequence.Address, missedInRowCount))

			if missedInRowCount == MissedInRowThreshold {
				newSystemEvent, err := t.newSystemEvent(validatorSequence, model.SystemEventMissedMInRow, systemEventRawData{
					"threshold": MissedInRowThreshold,
				})
				if err != nil {
					return nil, err
				}

				systemEvents = append(systemEvents, newSystemEvent)
			}
		}
	}
	return systemEvents, nil
}

// getTotalMissed get total missed count for given slice of validator sequences
func (t systemEventCreatorTask) getTotalMissed(validatorSequences []model.ValidatorSeq) int64 {
	var totalMissedCount int64 = 0
	for _, validatorSequence := range validatorSequences {
		if t.isNotValidated(validatorSequence) {
			totalMissedCount += 1
		}
	}

	return totalMissedCount
}

// getMissedInRow get number of validator sequences missed in the row
func (t systemEventCreatorTask) getMissedInRow(validatorSequences []model.ValidatorSeq, limit int64) int64 {
	if int64(len(validatorSequences)) > MissedInRowThreshold {
		validatorSequences = validatorSequences[:limit]
	}

	var missedInRowCount int64 = 0
	prevValidated := false
	for _, validatorSequence := range validatorSequences {
		if t.isNotValidated(validatorSequence) {
			if !prevValidated {
				missedInRowCount += 1
			}
			prevValidated = false
		} else {
			prevValidated = true
		}
	}

	return missedInRowCount
}

// isNotValidated check if validator has validated the block at height
func (t systemEventCreatorTask) isNotValidated(validatorSequence model.ValidatorSeq) bool {
	return validatorSequence.PrecommitValidated != nil && !*validatorSequence.PrecommitValidated
}

func (t systemEventCreatorTask) isValidated(validatorSequence model.ValidatorSeq) bool {
	return !t.isNotValidated(validatorSequence)
}

func (t *systemEventCreatorTask) getActiveSetPresenceChangeSystemEvents(currHeightValidatorSequences []model.ValidatorSeq, prevHeightValidatorSequences []model.ValidatorSeq) ([]*model.SystemEvent, error) {
	var systemEvents []*model.SystemEvent

	for _, currentValidatorSequence := range currHeightValidatorSequences {
		joined := true
		for _, prevValidatorSequence := range prevHeightValidatorSequences {
			if currentValidatorSequence.Address == prevValidatorSequence.Address {
				joined = false
				break
			}
		}

		if joined {
			logger.Debug(fmt.Sprintf("address %s joined active set", currentValidatorSequence.Address))

			newSystemEvent, err := t.newSystemEvent(currentValidatorSequence, model.SystemEventJoinedActiveSet, systemEventRawData{})
			if err != nil {
				return nil, err
			}

			systemEvents = append(systemEvents, newSystemEvent)
		}
	}

	for _, prevValidatorSequence := range prevHeightValidatorSequences {
		left := true
		for _, currentValidatorSequence := range currHeightValidatorSequences {
			if prevValidatorSequence.Address == currentValidatorSequence.Address {
				left = false
				break
			}
		}

		if left {
			logger.Debug(fmt.Sprintf("address %s joined active set", prevValidatorSequence.Address))

			newSystemEvent, err := t.newSystemEvent(prevValidatorSequence, model.SystemEventLeftActiveSet, systemEventRawData{})
			if err != nil {
				return nil, err
			}

			systemEvents = append(systemEvents, newSystemEvent)
		}
	}

	return systemEvents, nil
}

func (t *systemEventCreatorTask) getValueChangeSystemEvents(currHeightValidatorSequences []model.ValidatorSeq, prevHeightValidatorSequences []model.ValidatorSeq) ([]*model.SystemEvent, error) {
	var systemEvents []*model.SystemEvent
	for _, validatorSequence := range currHeightValidatorSequences {
		for _, prevValidatorSequence := range prevHeightValidatorSequences {
			if validatorSequence.Address == prevValidatorSequence.Address {
				newSystemEvent, err := t.getActiveEscrowBalanceChange(validatorSequence, prevValidatorSequence)
				if err != nil {
					if err != ErrActiveEscrowBalanceOutsideOfRange {
						return nil, err
					}
				} else {
					logger.Debug(fmt.Sprintf("active escrow balance change for address %s occured [kind=%s]", validatorSequence.Address, newSystemEvent.Kind))
					systemEvents = append(systemEvents, newSystemEvent)
				}

				newSystemEvent, err = t.getCommissionChange(validatorSequence, prevValidatorSequence)
				if err != nil {
					if err != ErrCommissionOutsideOfRange {
						return nil, err
					}
				} else {
					logger.Debug(fmt.Sprintf("commission change for address %s occured [kind=%s]", validatorSequence.Address, newSystemEvent.Kind))
					systemEvents = append(systemEvents, newSystemEvent)
				}
			}
		}
	}

	return systemEvents, nil
}

func (t *systemEventCreatorTask) getActiveEscrowBalanceChange(currValidatorSeq model.ValidatorSeq, prevValidatorSeq model.ValidatorSeq) (*model.SystemEvent, error) {
	currValue := currValidatorSeq.ActiveEscrowBalance.Int64()
	prevValue := prevValidatorSeq.ActiveEscrowBalance.Int64()
	roundedChangeRate := t.getRoundedChangeRate(currValue, prevValue)
	roundedAbsChangeRate := math.Abs(roundedChangeRate)

	var kind model.SystemEventKind
	if roundedAbsChangeRate >= 0.1 && roundedAbsChangeRate < 1 {
		kind = model.SystemEventActiveEscrowBalanceChange1
	} else if roundedAbsChangeRate >= 1 && roundedAbsChangeRate < 10 {
		kind = model.SystemEventActiveEscrowBalanceChange2
	} else if roundedAbsChangeRate >= 10 {
		kind = model.SystemEventActiveEscrowBalanceChange3
	} else {
		return nil, ErrActiveEscrowBalanceOutsideOfRange
	}

	return t.newSystemEvent(currValidatorSeq, kind, systemEventRawData{
		"before": prevValue,
		"after":  currValue,
		"change": roundedChangeRate,
	})
}

func (t *systemEventCreatorTask) getCommissionChange(currValidatorSeq model.ValidatorSeq, prevValidatorSeq model.ValidatorSeq) (*model.SystemEvent, error) {
	currValue := currValidatorSeq.Commission.Int64()
	prevValue := prevValidatorSeq.Commission.Int64()
	roundedChangeRate := t.getRoundedChangeRate(currValue, prevValue)
	roundedAbsChangeRate := math.Abs(roundedChangeRate)

	var kind model.SystemEventKind
	if roundedAbsChangeRate >= 0.1 && roundedAbsChangeRate < 1 {
		kind = model.SystemEventCommissionChange1
	} else if roundedAbsChangeRate >= 1 && roundedAbsChangeRate < 10 {
		kind = model.SystemEventCommissionChange2
	} else if roundedAbsChangeRate >= 10 {
		kind = model.SystemEventCommissionChange3
	} else {
		return nil, ErrCommissionOutsideOfRange
	}

	return t.newSystemEvent(currValidatorSeq, kind, systemEventRawData{
		"before": prevValue,
		"after":  currValue,
		"change": roundedChangeRate,
	})
}

func (t *systemEventCreatorTask) getRoundedChangeRate(currValue int64, prevValue int64) float64 {
	var changeRate float64

	if prevValue == 0 {
		changeRate = float64(currValue)
	} else {
		changeRate = (float64(1) - (float64(currValue) / float64(prevValue))) * 100
	}

	roundedChangeRate := math.Round(changeRate/0.1) * 0.1
	return roundedChangeRate
}

func (t *systemEventCreatorTask) newSystemEvent(seq model.ValidatorSeq, kind model.SystemEventKind, data map[string]interface{}) (*model.SystemEvent, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &model.SystemEvent{
		Height: seq.Height,
		Time:   seq.Time,
		Actor:  seq.Address,
		Kind:   kind,
		Data:   types.Jsonb{RawMessage: marshaledData},
	}, nil
}
