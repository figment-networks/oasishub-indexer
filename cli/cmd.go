package cli

import (
	"context"
	"fmt"

	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

func runCmd(cfg *config.Config, flags Flags) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	client, err := initClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	cmdHandlers := usecase.NewCmdHandlers(cfg, db, client)

	logger.Info(fmt.Sprintf("executing cmd %s ...", flags.runCommand), logger.Field("app", "cli"))

	ctx := context.Background()
	switch flags.runCommand {
	case "status":
		cmdHandlers.GetStatus.Handle(ctx)
	case "indexer:index":
		cmdHandlers.IndexerIndex.Handle(ctx, flags.batchSize)
	case "indexer:backfill":
		cmdHandlers.IndexerBackfill.Handle(ctx, flags.parallel, flags.force)
	case "indexer:reindex":
		cmdHandlers.IndexerReindex.Handle(ctx, flags.parallel, flags.force, flags.startReindexHeight, flags.endReindexHeight, flags.targetIds)
	case "indexer:summarize":
		cmdHandlers.IndexerSummarize.Handle(ctx)
	case "indexer:purge":
		cmdHandlers.IndexerPurge.Handle(ctx)
	case "validators:decorate":
		cmdHandlers.DecorateValidators.Handle(ctx, flags.filePath)
	default:
		return errors.New(fmt.Sprintf("command %s not found", flags.runCommand))
	}
	return nil
}
