package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	TaskNameHeightMetaRetriever = "HeightMetaRetriever"
)

func NewHeightMetaRetrieverTask(c client.ChainClient) pipeline.Task {
	return &heightMetaRetrieverTask{
		client: c,
	}
}

type heightMetaRetrieverTask struct {
	client client.ChainClient
}

type HeightMeta struct {
	Height       int64
	Time         types.Time
	AppVersion   uint64
	BlockVersion uint64
}

func (t *heightMetaRetrieverTask) GetName() string {
	return TaskNameHeightMetaRetriever
}

func (t *heightMetaRetrieverTask) Run(ctx context.Context, p pipeline.Payload) error {

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

	meta, err := t.client.GetMetaByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	payload.HeightMeta = HeightMeta{
		Height:       meta.GetHeight(),
		Time:         *types.NewTimeFromTimestamp(*meta.GetTime()),
		AppVersion:   meta.GetAppVersion(),
		BlockVersion: meta.GetBlockVersion(),
	}
	return nil
}
