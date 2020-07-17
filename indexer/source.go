package indexer

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

var (
	_ pipeline.Source = (*source)(nil)

	ErrNothingToProcess = errors.New("nothing to process")
)

type sourceStore interface {
	FindMostRecent() (*model.Syncable, error)
}

func NewSource(cfg *config.Config, db sourceStore, client client.ChainClient, versionNumber int64, batchSize int64) *source {
	src := &source{
		cfg:    cfg,
		db:     db,
		client: client,

		versionNumber: versionNumber,
		batchSize:     batchSize,
	}
	return src.init()
}

type source struct {
	cfg           *config.Config
	db            sourceStore
	client        client.ChainClient
	versionNumber int64
	batchSize     int64

	indexVersion  int64
	currentHeight int64
	startHeight   int64
	endHeight     int64
	err           error
}

func (s *source) Next(context.Context, pipeline.Payload) bool {
	if s.err == nil && s.currentHeight < s.endHeight {
		s.currentHeight = s.currentHeight + 1
		return true
	}
	return false
}

func (s *source) Current() int64 {
	return s.currentHeight
}

func (s *source) Err() error {
	return s.err
}

func (s *source) Len() int64 {
	return s.endHeight - s.startHeight + 1
}

func (s *source) init() *source {
	err := s.setStartHeight()
	if err != nil {
		s.err = err
		return s
	}
	err = s.setEndHeight()
	if err != nil {
		s.err = err
		return s
	}
	if err := s.validate(); err != nil {
		s.err = err
		return s
	}
	return s
}

func (s *source) setStartHeight() error {
	var startH int64
	syncable, err := s.db.FindMostRecent()
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
	return nil
}

func (s *source) setEndHeight() error {
	syncableFromNode, err := s.client.GetHead()
	if err != nil {
		return err
	}
	endH := syncableFromNode.GetChain().GetHeight()

	if s.batchSize > 0 && endH-s.startHeight > s.batchSize {
		endOfBatch := (s.startHeight + s.batchSize) - 1
		endH = endOfBatch
	}
	s.endHeight = endH
	return nil
}

func (s *source) validate() error {
	blocksToSyncCount := s.endHeight - s.startHeight
	if blocksToSyncCount == 0 {
		return ErrNothingToProcess
	}
	return nil
}
