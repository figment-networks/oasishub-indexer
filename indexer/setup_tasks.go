package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/metric"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	HeightMetaRetrieverTaskName = "HeightMetaRetriever"
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

func (t *heightMetaRetrieverTask) GetName() string {
	return HeightMetaRetrieverTaskName
}

func (t *heightMetaRetrieverTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

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
