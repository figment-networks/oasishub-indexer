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
	"strconv"
	"strings"
)

type Flags struct {
	configPath  string
	runCommand  string
	showVersion bool

	batchSize  int64
	parallel   bool
	force      bool
	versionIds versionIds
	targetIds  targetIds
}

// versionIds holds comma-separated list of version ids
type versionIds []int64

func (i *versionIds) String() string {
	return fmt.Sprint(*i)
}

func (i *versionIds) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("versionIds flag already set")
	}
	*i = splitIds(value)
	return nil
}

// targetIds holds comma-separated list of target ids
type targetIds []int64

func (i *targetIds) String() string {
	return fmt.Sprint(*i)
}

func (i *targetIds) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("targetIds flag already set")
	}
	*i = splitIds(value)
	return nil
}

func splitIds(ids string) []int64 {
	var parsedIds []int64
	for _, rawId := range strings.Split(ids, ",") {
		id, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			panic(err)
		}
		parsedIds = append(parsedIds, id)
	}
	return parsedIds
}

func (c *Flags) Setup() {
	flag.BoolVar(&c.showVersion, "v", false, "Show application version")
	flag.StringVar(&c.configPath, "config", "", "Path to config")
	flag.StringVar(&c.runCommand, "cmd", "", "Command to run")

	flag.Int64Var(&c.batchSize, "batch_size", 0, "pipeline batch size")
	flag.BoolVar(&c.parallel, "parallel", false, "should backfill be run in parallel with indexing")
	flag.BoolVar(&c.force, "force", false, "remove existing reindexing reports")
	flag.Var(&c.targetIds, "target_ids", "comma separated list of integers")
}

// Run executes the command line interface
func Run() {
	flags := Flags{}
	flags.Setup()
	flag.Parse()

	if flags.showVersion {
		fmt.Println(config.VersionString())
		return
	}

	// Initialize configuration
	cfg, err := initConfig(flags.configPath)
	if err != nil {
		panic(fmt.Errorf("error initializing config [ERR: %+v]", err))
	}

	// Initialize logger
	if err = initLogger(cfg); err != nil {
		panic(fmt.Errorf("error initializing logger [ERR: %+v]", err))
	}

	// Initialize error reporting
	initErrorReporting(cfg)

	if flags.runCommand == "" {
		terminate(errors.New("command is required"))
	}

	if err := startCommand(cfg, flags); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, flags Flags) error {
	switch flags.runCommand {
	case "migrate":
		return startMigrations(cfg)
	case "server":
		return startServer(cfg)
	case "worker":
		return startWorker(cfg)
	default:
		return runCmd(cfg, flags)
	}
}

func terminate(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func initConfig(path string) (*config.Config, error) {
	cfg := config.New()

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	if path != "" {
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
	return logger.Init(cfg)
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
