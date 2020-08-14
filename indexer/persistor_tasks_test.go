package indexer

import (
	"context"
	"fmt"
	"testing"
	"time"

	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
)

func TestSyncerPersistor_Run(t *testing.T) {
	sync := &model.Syncable{
		Height: 20,
		Time:   *types.NewTimeFromTime(time.Now()),
	}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with syncable", nil},
		{"returns error if database errors", errTestDbCreate},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockSyncerPersistorTaskStore(ctrl)

			task := NewSyncerPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight: 20,
				Syncable:      sync,
			}

			dbMock.EXPECT().CreateOrUpdate(sync).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}

func TestBlockSeqPersistor_Run(t *testing.T) {
	seq := &model.BlockSeq{
		Sequence: &model.Sequence{
			Height: 20,
			Time:   *types.NewTimeFromTime(time.Date(1987, 12, 11, 14, 0, 0, 0, time.UTC)),
		},
		TransactionsCount: 10,
	}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with block sequence", nil},
		{"returns error if database errors", errTestDbCreate},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[new] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockBlockSeqPersistorTaskStore(ctrl)

			task := NewBlockSeqPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:    20,
				NewBlockSequence: seq,
			}

			dbMock.EXPECT().Create(seq).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})

		t.Run(fmt.Sprintf("[updated] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockBlockSeqPersistorTaskStore(ctrl)

			task := NewBlockSeqPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:        20,
				UpdatedBlockSequence: seq,
			}

			dbMock.EXPECT().Save(seq).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}

func TestValidatorSeqPersistor_Run(t *testing.T) {
	newValidatorSeq := func() model.ValidatorSeq {
		return model.ValidatorSeq{
			Sequence: &model.Sequence{
				Height: 20,
				Time:   *types.NewTimeFromTime(time.Date(1987, 12, 11, 14, 0, 0, 0, time.UTC)),
			},
			EntityUID: randString(5),
		}
	}

	seq := []model.ValidatorSeq{newValidatorSeq(), newValidatorSeq()}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with all validator sequences", nil},
		{"returns error if database errors", errTestDbCreate},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[new] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorSeqPersistorTaskStore(ctrl)

			task := NewValidatorSeqPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:         20,
				NewValidatorSequences: seq,
			}

			for _, s := range seq {
				createSeq := s
				dbMock.EXPECT().Create(&createSeq).Return(tt.expectErr).Times(1)
				if tt.expectErr != nil {
					// don't expect any more calls
					break
				}
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})

		t.Run(fmt.Sprintf("[updated] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorSeqPersistorTaskStore(ctrl)

			task := NewValidatorSeqPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:             20,
				UpdatedValidatorSequences: seq,
			}

			for _, s := range seq {
				saveSeq := s
				dbMock.EXPECT().Save(&saveSeq).Return(tt.expectErr).Times(1)
				if tt.expectErr != nil {
					// don't expect any more calls
					break
				}
			}
			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}

func TestValidatorAggPersistor_Run(t *testing.T) {
	newValidatorAgg := func() model.ValidatorAgg {
		return model.ValidatorAgg{
			Aggregate: &model.Aggregate{
				StartedAtHeight: 20,
				StartedAt:       *types.NewTimeFromTime(time.Date(1988, 12, 11, 14, 0, 0, 0, time.UTC)),
			},
			EntityUID: randString(5),
		}
	}

	seq := []model.ValidatorAgg{newValidatorAgg(), newValidatorAgg()}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with all validator aggregates", nil},
		{"returns error if database errors", errTestDbCreate},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[new] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorAggPersistorTaskStore(ctrl)

			task := NewValidatorAggPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:           20,
				NewAggregatedValidators: seq,
			}

			for _, s := range seq {
				createSeq := s
				dbMock.EXPECT().Create(&createSeq).Return(tt.expectErr).Times(1)
				if tt.expectErr != nil {
					// don't expect any more calls
					break
				}
			}

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})

		t.Run(fmt.Sprintf("[updated] %v", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorAggPersistorTaskStore(ctrl)

			task := NewValidatorAggPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:               20,
				UpdatedAggregatedValidators: seq,
			}

			for _, s := range seq {
				saveSeq := s
				dbMock.EXPECT().Save(&saveSeq).Return(tt.expectErr).Times(1)
				if tt.expectErr != nil {
					// don't expect any more calls
					break
				}
			}
			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}
