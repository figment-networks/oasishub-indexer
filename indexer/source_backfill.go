package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

var (
	_ pipeline.Source = (*backfillSource)(nil)

	ErrNothingToBackfill = errors.New("nothing to backfill")
)

type BackfillSourceStore interface {
	FindFirstByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error)
	FindMostRecentByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error)
}

func NewBackfillSource(cfg *config.Config, db BackfillSourceStore, civ int64) (*backfillSource, error) {
	src := &backfillSource{
		cfg: cfg,
		db:  db,

		currentIndexVersion: civ,
	}

	if err := src.init(); err != nil {
		return nil, err
	}

	return src, nil
}

type backfillSource struct {
	cfg *config.Config
	db  BackfillSourceStore

	currentIndexVersion int64

	startHeight   int64
	currentHeight int64
	endHeight     int64
	err           error
}

func (s *backfillSource) Next(context.Context, pipeline.Payload) bool {
	if s.err == nil && s.currentHeight < s.endHeight {
		s.currentHeight = s.currentHeight + 1
		return true
	}
	return false
}

func (s *backfillSource) Current() int64 {
	return s.currentHeight
}

func (s *backfillSource) Err() error {
	return s.err
}

func (s *backfillSource) Len() int64 {
	return s.endHeight - s.startHeight + 1
}

func (s *backfillSource) init() error {
	if err := s.setStartHeight(); err != nil {
		return err
	}
	if err := s.setEndHeight(); err != nil {
		return err
	}
	return nil
}

func (s *backfillSource) setStartHeight() error {
	syncable, err := s.db.FindFirstByDifferentIndexVersion(s.currentIndexVersion)
	if err != nil {
		if err == store.ErrNotFound {
			return errors.Wrap(ErrNothingToBackfill, fmt.Sprintf("[currentIndexVersion=%d] Reason: everything is up to date", s.currentIndexVersion))
		}
		return err
	}

	s.currentHeight = syncable.Height
	s.startHeight = syncable.Height
	return nil
}

func (s *backfillSource) setEndHeight() error {
	syncable, err := s.db.FindMostRecentByDifferentIndexVersion(s.currentIndexVersion)
	if err != nil {
		return err
	}

	s.endHeight = syncable.Height
	return nil
}
