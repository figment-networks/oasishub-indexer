package indexing

import (
	"context"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/types"
)

func NewHeightMetaRetrieverTask(c *client.Client) pipeline.Task {
	return &heightMetaRetrieverTask{
		client: c,
	}
}

type heightMetaRetrieverTask struct {
	client *client.Client
}

type HeightMeta struct {
	Height       int64
	Time         types.Time
	AppVersion   uint64
	BlockVersion uint64
}

func (t *heightMetaRetrieverTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)
	block, err := t.client.Block.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	payload.HeightMeta = HeightMeta{
		Height:       block.GetBlock().GetHeader().GetHeight(),
		Time:         *types.NewTimeFromTimestamp(*block.GetBlock().GetHeader().GetTime()),
		AppVersion:   block.GetBlock().GetHeader().GetVersion().GetApp(),
		BlockVersion: block.GetBlock().GetHeader().GetVersion().GetBlock(),
	}
	return nil
}
