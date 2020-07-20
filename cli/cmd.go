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

func runCmd(cfg *config.Config, cmdName string, filePath string) error {
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

	logger.Info(fmt.Sprintf("executing cmd %s ...", cmdName), logger.Field("app", "cli"))

	switch cmdName {
	case "run_indexer":
		ctx := context.Background()
		cmdHandlers.RunIndexer.Handle(ctx)
	case "purge_indexer":
		ctx := context.Background()
		cmdHandlers.PurgeIndexer.Handle(ctx)
	case "summarize_indexer":
		ctx := context.Background()
		cmdHandlers.SummarizeIndexer.Handle(ctx)
	case "decorate_validators":
		ctx := context.Background()
		ctxWithFilePath := context.WithValue(ctx, validator.CtxFilePath, filePath)
		cmdHandlers.DecorateValidators.Handle(ctxWithFilePath)
	default:
		return errors.New(fmt.Sprintf("command %s not found", cmdName))
	}
	return nil
}
