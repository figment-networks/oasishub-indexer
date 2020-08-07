package indexer

import (
	"context"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

var (
	_ pipeline.Source = (*indexSource)(nil)

	ErrNothingToProcess = errors.New("nothing to process")
)

func NewIndexSource(cfg *config.Config, db *store.Store, client *client.Client, startHeight int64, batchSize int64) (*indexSource, error) {
	src := &indexSource{
		cfg:    cfg,
		db:     db,
		client: client,

		batchSize:     batchSize,

		startHeight:   startHeight,
		currentHeight: startHeight,
	}
	if err := src.init(); err != nil {
		return nil, err
	}
	return src, nil
}

type indexSource struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	batchSize int64

	startHeight   int64
	currentHeight int64
	endHeight     int64
	err           error
}

func (s *indexSource) Next(context.Context, pipeline.Payload) bool {
	if s.err == nil && s.currentHeight < s.endHeight {
		s.currentHeight = s.currentHeight + 1
		return true
	}
	return false
}

func (s *indexSource) Current() int64 {
	return s.currentHeight
}

func (s *indexSource) Err() error {
	return s.err
}

func (s *indexSource) Len() int64 {
	return s.endHeight - s.startHeight + 1
}

func (s *indexSource) init() error {
	if err := s.setStartHeight(); err != nil {
		return err
	}
	if err := s.setEndHeight(); err != nil {
		return err
	}
	if err := s.validate(); err != nil {
		return err
	}
	return nil
}

func (s *indexSource) setStartHeight() error {
	if s.startHeight == 0 {
		var startH int64
		syncable, err := s.db.Syncables.FindMostRecent()
		if err != nil {
			if err != store.ErrNotFound {
				return err
			}
			// No syncables found, get first block number from config
			startH = s.cfg.FirstBlockHeight
		} else {
			// Reindex if last syncable failed
			if syncable.ProcessedAt == nil {
				startH = syncable.Height
			} else {
				startH = syncable.Height + 1
			}
		}
		s.currentHeight = startH
		s.startHeight = startH
	}
	return nil
}

func (s *indexSource) setEndHeight() error {
	syncableFromNode, err := s.client.Chain.GetHead()
	if err != nil {
		return err
	}
	endH := syncableFromNode.GetHeight()

	if s.batchSize > 0 && endH-s.startHeight > s.batchSize {
		endOfBatch := (s.startHeight + s.batchSize) - 1
		endH = endOfBatch
	}

	s.endHeight = endH
	return nil
}

func (s *indexSource) validate() error {
	blocksToSyncCount := s.endHeight - s.startHeight
	if blocksToSyncCount == 0 && s.batchSize != 1 {
		return ErrNothingToProcess
	}
	return nil
}
