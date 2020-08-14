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

	aggs, err := uc.db.ValidatorAgg.GetAllForHeightGreaterThan(*height)
	if err != nil {
		return nil, err
	}

	seqs, err := uc.db.ValidatorSeq.FindByHeight(*height)
	if len(seqs) == 0 || err != nil {
		indexingPipeline, err := indexer.NewPipeline(uc.cfg, uc.db, uc.client)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		payload, err := indexingPipeline.Run(ctx, indexer.RunConfig{
			Height:           *height,
			DesiredTargetIDs: []int64{indexer.IndexTargetValidatorSequences},
			Dry:              true,
		})
		if err != nil {
			return nil, err
		}

		seqs = payload.NewValidatorSequences
	}

	return ToSeqListView(seqs, aggs), nil
}
