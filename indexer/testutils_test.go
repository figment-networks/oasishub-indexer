package indexer

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

func setup(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	err := logger.InitTestLogger()
	if err != nil {
		t.Fatal(err)
	}
}

func ctxWithReport(modelID types.ID) context.Context {
	ctx := context.Background()
	report := &model.Report{
		Model: &model.Model{ID: modelID},
	}

	return context.WithValue(ctx, CtxReport, report)
}
