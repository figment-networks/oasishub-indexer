package indexer

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/model"
)

var (
	_ pipeline.Source = (*reindexSource)(nil)
)

type ReindexSourceStore interface {
	FindMostRecent() (*model.Syncable, error)
}

func NewReindexSource(cfg *config.Config, db ReindexSourceStore, startHeight, endHeight int64) (*reindexSource, error) {
	src := &reindexSource{
		cfg: cfg,
		db:  db,
	}

	if err := src.init(startHeight, endHeight); err != nil {
		return nil, err
	}

	return src, nil
}

type reindexSource struct {
	cfg *config.Config
	db  ReindexSourceStore

	startHeight   int64
	currentHeight int64
	endHeight     int64
	err           error
}

func (s *reindexSource) Next(context.Context, pipeline.Payload) bool {
	if s.err == nil && s.currentHeight < s.endHeight {
		s.currentHeight = s.currentHeight + 1
		return true
	}
	return false
}

func (s *reindexSource) Skip(stageName pipeline.StageName) bool {
	return false
}

func (s *reindexSource) Current() int64 {
	return s.currentHeight
}

func (s *reindexSource) Err() error {
	return s.err
}

func (s *reindexSource) Len() int64 {
	return s.endHeight - s.startHeight + 1
}

func (s *reindexSource) init(startHeight, endHeight int64) error {
	if err := s.setStartHeight(startHeight); err != nil {
		return err
	}
	if err := s.setEndHeight(endHeight); err != nil {
		return err
	}
	return nil
}

func (s *reindexSource) setStartHeight(startHeight int64) error {
	startH := startHeight
	if startH < 1 {
		startH = 1
	}

	s.currentHeight = startH
	s.startHeight = startH
	return nil
}

func (s *reindexSource) setEndHeight(endHeight int64) error {
	if endHeight > 0 {
		s.endHeight = endHeight
		return nil
	}

	syncable, err := s.db.FindMostRecent()
	if err != nil {
		return err
	}

	s.endHeight = syncable.Height
	return nil
}
