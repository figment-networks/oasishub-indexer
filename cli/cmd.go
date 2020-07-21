package cli

import (
	"context"
	"fmt"

	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/usecase/validator"
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
	case "indexer_start":
		cmdHandlers.StartIndexer.Handle(ctx, flags.batchSize)
	case "indexer_backfill":
		cmdHandlers.BackfillIndexer.Handle(ctx, flags.parallel, flags.force, flags.targetIds)
	case "indexer_summarize":
		cmdHandlers.SummarizeIndexer.Handle(ctx)
	case "decorate_validators":
		ctx := context.Background()
		ctxWithFilePath := context.WithValue(ctx, validator.CtxFilePath, filePath)
		cmdHandlers.DecorateValidators.Handle(ctxWithFilePath)
	case "indexer_purge":
		cmdHandlers.PurgeIndexer.Handle(ctx)
	default:
		return errors.New(fmt.Sprintf("command %s not found", flags.runCommand))
	}
	return nil
}
