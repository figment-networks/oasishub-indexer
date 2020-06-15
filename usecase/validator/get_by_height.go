package validator

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/indexer"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client
}

func NewGetByHeightUseCase(cfg *config.Config, db *store.Store, client *client.Client) *getByHeightUseCase {
	return &getByHeightUseCase{
		cfg:    cfg,
		db:     db,
		client: client,
	}
}

func (uc *getByHeightUseCase) Execute(height *int64) (*SeqListView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}
	lastH := mostRecentSynced.Height

	// Show last synced height, if not provided
	if height == nil {
		height = &lastH
	}

	if *height > lastH {
		return nil, errors.New("height is not indexed yet")
	}

	models, err := uc.db.ValidatorSeq.FindByHeight(*height)
	if len(models) == 0 || err != nil {
		indexingPipeline, err := indexer.NewPipeline(uc.cfg, uc.db, uc.client)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		payload, err := indexingPipeline.Run(ctx, indexer.RunConfig{
			Height: *height,
			DesiredTargetID: indexer.TargetIndexValidatorSequences,
			Dry:    true,
		})
		if err != nil {
			return nil, err
		}

		models = payload.NewValidatorSequences
	}

	return ToSeqListView(models), nil
}
