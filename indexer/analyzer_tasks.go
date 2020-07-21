package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
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

	MaxValidatorSequences int64 = 1000
	MissedForMaxThreshold int64 = 50
	MissedInRowThreshold  int64 = 50
)

// NewSystemEventCreatorTask creates system events
func NewSystemEventCreatorTask(vStore ValidatorSeqStore) *systemEventCreatorTask {
	return &systemEventCreatorTask{
		ValidatorSeqStore: vStore,
	}
}

type ValidatorSeqStore interface {
	FindByHeight(int64) ([]model.ValidatorSeq, error)
	FindLastByAddress(string, int64) ([]model.ValidatorSeq, error)
}

type systemEventCreatorTask struct {
	ValidatorSeqStore
}

func (t *systemEventCreatorTask) GetName() string {
	return SystemEventCreatorTaskName
}

func (t *systemEventCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", "Analyzer", t.GetName(), payload.CurrentHeight))

	currHeightValidatorSequences := append(payload.NewValidatorSequences, payload.UpdatedValidatorSequences...)
	prevHeightValidatorSequences, err := t.ValidatorSeqStore.FindByHeight(payload.CurrentHeight - 1)
	if err != nil {
		if err != store.ErrNotFound {
			return err
		}
	}

	activeEscrowBalanceChangeSystemEvents := t.getActiveEscrowBalanceChangeSystemEvents(currHeightValidatorSequences, prevHeightValidatorSequences)
	payload.SystemEvents = append(payload.SystemEvents, activeEscrowBalanceChangeSystemEvents...)

	activeSetPresenceChangeSystemEvents := t.getActiveSetPresenceChangeSystemEvents(currHeightValidatorSequences, prevHeightValidatorSequences)
	payload.SystemEvents = append(payload.SystemEvents, activeSetPresenceChangeSystemEvents...)

	missedBlocksSystemEvents, err := t.getMissedBlocksSystemEvents(currHeightValidatorSequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, missedBlocksSystemEvents...)

	return nil
}

func (t *systemEventCreatorTask) getMissedBlocksSystemEvents(currHeightValidatorSequences []model.ValidatorSeq) ([]*model.SystemEvent, error) {
	var systemEvents []*model.SystemEvent
	for _, validatorSequence := range currHeightValidatorSequences {
		// When current height validator has validated the block no need to check last records
		if t.isValidated(validatorSequence) {
			return systemEvents, nil
		}

		lastValidatorSequencesForAddress, err := t.ValidatorSeqStore.FindLastByAddress(validatorSequence.Address, MaxValidatorSequences)
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
				systemEvents = append(systemEvents, t.newSystemEventWithBody(validatorSequence, model.SystemEventMissedMofN))
			}

			missedInRowCount := t.getMissedInRow(validatorSequencesToCheck, MissedInRowThreshold)

			logger.Debug(fmt.Sprintf("total missed blocks in a row for address %s: %d", validatorSequence.Address, missedInRowCount))

			if missedInRowCount == MissedInRowThreshold {
				systemEvents = append(systemEvents, t.newSystemEventWithBody(validatorSequence, model.SystemEventMissedMInRow))
			}
		}
	}
	return systemEvents, nil
}

func (t systemEventCreatorTask) getTotalMissed(validatorSequences []model.ValidatorSeq) int64 {
	var totalMissedCount int64 = 0
	for _, validatorSequence := range validatorSequences {
		if t.isNotValidated(validatorSequence) {
			totalMissedCount += 1
		}
	}

	return totalMissedCount
}

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

func (t systemEventCreatorTask) isNotValidated(validatorSequence model.ValidatorSeq) bool {
	return validatorSequence.PrecommitValidated != nil && !*validatorSequence.PrecommitValidated
}

func (t systemEventCreatorTask) isValidated(validatorSequence model.ValidatorSeq) bool {
	return !t.isNotValidated(validatorSequence)
}

func (t *systemEventCreatorTask) getActiveSetPresenceChangeSystemEvents(currHeightValidatorSequences []model.ValidatorSeq, prevHeightValidatorSequences []model.ValidatorSeq) []*model.SystemEvent {
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

			systemEvents = append(systemEvents, t.newSystemEventWithBody(currentValidatorSequence, model.SystemEventJoinedActiveSet))
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

			systemEvents = append(systemEvents, t.newSystemEventWithBody(prevValidatorSequence, model.SystemEventLeftActiveSet))
		}
	}

	return systemEvents
}

func (t *systemEventCreatorTask) getActiveEscrowBalanceChangeSystemEvents(currHeightValidatorSequences []model.ValidatorSeq, prevHeightValidatorSequences []model.ValidatorSeq) []*model.SystemEvent {
	var systemEvents []*model.SystemEvent
	for _, validatorSequence := range currHeightValidatorSequences {
		for _, prevValidatorSequence := range prevHeightValidatorSequences {
			if validatorSequence.Address == prevValidatorSequence.Address {
				changeRate := (float64(1) - (float64(validatorSequence.ActiveEscrowBalance.Int64()) / float64(prevValidatorSequence.ActiveEscrowBalance.Int64()))) * 100

				kind, err := t.getActiveEscrowBalanceChangeKind(changeRate)
				if err == nil {
					logger.Debug(fmt.Sprintf("active escrow balance change for address %s occured [kind=%s]", validatorSequence.Address, kind))

					systemEvents = append(systemEvents, t.newSystemEventWithBody(validatorSequence, *kind))
				}
			}
		}
	}

	return systemEvents
}

func (t *systemEventCreatorTask) getActiveEscrowBalanceChangeKind(changeRate float64) (*model.SystemEventKind, error) {
	roundedAbsChangeRate := math.Round(math.Abs(changeRate) / 0.1) * 0.1

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

	return &kind, nil
}

func (t *systemEventCreatorTask) newSystemEventWithBody(seq model.ValidatorSeq, kind model.SystemEventKind) *model.SystemEvent {
	return &model.SystemEvent{
		Height: seq.Height,
		Time:   seq.Time,
		Actor:  seq.Address,
		Kind:   kind,
	}
}
