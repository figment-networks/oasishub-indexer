package indexing

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

var (
	ErrRunningSequentialReindex = errors.New("indexing skipped because sequential reindex hasn't finished yet")
)

type indexUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client
}

func NewIndexUseCase(cfg *config.Config, db *store.Store, c *client.Client) *indexUseCase {
	return &indexUseCase{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (uc *indexUseCase) Execute(ctx context.Context, batchSize int64) error {
	if err := uc.canExecute(); err != nil {
		return err
	}

	indexingPipeline := indexer.NewPipeline(uc.cfg, uc.db, uc.client)

	return indexingPipeline.Index(ctx, indexer.IndexCfg{
		BatchSize: batchSize,
	})
}

// canExecute checks if sequential reindex is already running
// if is it running we skip indexing
func (uc *indexUseCase) canExecute() error {
	if _, err := uc.db.Reports.FindNotCompletedByKind(model.ReportKindSequentialReindex); err != nil {
		if err == store.ErrNotFound {
			return nil
		}
		return err
	}
	return ErrRunningSequentialReindex
}

