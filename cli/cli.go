package cli

import (
	"flag"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/figment-networks/oasishub-indexer/utils/reporting"
	"github.com/pkg/errors"
)

// Run executes the command line interface
func Run() {
	defer reporting.RecoverError()

	var configPath string
	var runCommand string
	var showVersion bool

	flag.BoolVar(&showVersion, "v", false, "Show application version")
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.StringVar(&runCommand, "cmd", "", "Command to run")
	flag.Parse()

	if showVersion {
		fmt.Println(versionString())
		return
	}

	cfg, err := initConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("error initializing config [ERR: %+v]", err))
	}

	if err = initLogger(cfg); err != nil {
		panic(fmt.Errorf("error initializing logger [ERR: %+v]", err))
	}

	initErrorReporting(cfg)

	if runCommand == "" {
		terminate(errors.New("command is required"))
	}

	if err := startCommand(cfg, runCommand); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, name string) error {
	switch name {
	case "migrate":
		return startMigrations(cfg)
	case "server":
		return startServer(cfg)
	case "worker":
		return startWorker(cfg)
	default:
		return runCmd(cfg, name)
	}
}

func terminate(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func initConfig(path string) (*config.Config, error) {
	cfg := config.New()

	if path == "" {
		if err := config.FromEnv(cfg); err != nil {
			return nil, err
		}
	} else {
		if err := config.FromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initLogger(cfg *config.Config) error {
	_, err := logger.Init(cfg)
	return err
}

func initClient(cfg *config.Config) (*client.Client, error) {
	return client.New(cfg.ProxyUrl)
}

func initStore(cfg *config.Config) (*store.Store, error) {
	db, err := store.New(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	db.SetDebugMode(cfg.Debug)

	return db, nil
}

func initErrorReporting(cfg *config.Config) {
	reporting.Init(cfg)
}
