package config

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/spf13/viper"
)

const (
	appPort       = "PORT"
	nodeUrl       = "NODE_URL"
	logLevel      = "LOG_LEVEL"
	logOutput     = "LOG_OUTPUT"
	goEnvironment = "GO_ENVIRONMENT"

	// Cli
	batchSize = "CLI_BATCH_SIZE_ARG"

	// Job
	pingInterval      = "PING_INTERVAL"
	syncInterval      = "SYNC_INTERVAL"
	processInterval   = "PROCESS_INTERVAL"
	cleanupInterval   = "CLEANUP_INTERVAL"
	pipelineBatchSize = "BLOCK_SYNC_BATCH_SIZE"
	firstBlockHeight  = "FIRST_BLOCK_HEIGHT"
	cleanupThreshold  = "CLEANUP_THRESHOLD"

	// Database
	dbUser        = "DB_USER"
	dbPassword    = "DB_PASSWORD"
	dbHost        = "DB_HOST"
	dbName        = "DB_NAME"
	dbDetailedLog = "DB_DETAILED_LOG"
	dbSSLMode	  = "DB_SSL_MODE"

	production  = "production"
	development = "development"
)

var (
	defaultAppPort       = "8081"
	defaultNodeUrl       = "http://localhost:8080"
	defaultLogLevel      = "info"
	defaultLogOutput     = "stdout"
	defaultGoEnvironment = development

	defaultBatchSize = "batchSize"

	defaultPipelineBatchSize int64 = 2
	defaultPingInterval            = "@every 20s"
	defaultSyncInterval            = "@every 20s"
	defaultProcessInterval         = "@every 20s"
	defaultCleanupInterval         = "@every 20s"
	defaultFirstBlockHeight  int64 = 1
	defaultCleanupThreshold  int64 = 1000

	defaultDbUser        = "postgres"
	defaultDbPassword    = "password"
	defaultDbHost        = "db"
	defaultDbName        = "app"
	defaultSSLMode       = "disable"
	dbDefaultDetailedLog = false

)

func init() {
	viper.SetDefault(appPort, defaultAppPort)
	viper.SetDefault(nodeUrl, defaultNodeUrl)
	viper.SetDefault(logLevel, defaultLogLevel)
	viper.SetDefault(logOutput, defaultLogOutput)
	viper.SetDefault(goEnvironment, defaultGoEnvironment)

	viper.SetDefault(batchSize, defaultBatchSize)

	viper.SetDefault(pipelineBatchSize, defaultPipelineBatchSize)
	viper.SetDefault(pingInterval, defaultPingInterval)
	viper.SetDefault(syncInterval, defaultSyncInterval)
	viper.SetDefault(processInterval, defaultProcessInterval)
	viper.SetDefault(cleanupInterval, defaultCleanupInterval)
	viper.SetDefault(firstBlockHeight, defaultFirstBlockHeight)
	viper.SetDefault(cleanupThreshold, defaultCleanupThreshold)

	viper.SetDefault(dbUser, defaultDbUser)
	viper.SetDefault(dbPassword, defaultDbPassword)
	viper.SetDefault(dbHost, defaultDbHost)
	viper.SetDefault(dbName, defaultDbName)
	viper.SetDefault(dbSSLMode, defaultSSLMode)
	viper.SetDefault(dbDetailedLog, dbDefaultDetailedLog)

	viper.BindEnv(appPort)
	viper.BindEnv(nodeUrl)
	viper.BindEnv(logLevel)
	viper.BindEnv(logOutput)
	viper.BindEnv(goEnvironment)

	viper.BindEnv(batchSize)

	viper.BindEnv(pipelineBatchSize)
	viper.BindEnv(pingInterval)
	viper.BindEnv(syncInterval)
	viper.BindEnv(processInterval)
	viper.BindEnv(cleanupInterval)
	viper.BindEnv(firstBlockHeight)
	viper.BindEnv(cleanupThreshold)

	viper.BindEnv(dbUser)
	viper.BindEnv(dbPassword)
	viper.BindEnv(dbHost)
	viper.BindEnv(dbName)
	viper.BindEnv(dbSSLMode)
	viper.BindEnv(dbDetailedLog)
}

func AppPort() string {
	return viper.GetString(appPort)
}

func NodeUrl() string {
	return viper.GetString(nodeUrl)
}

func LogLevel() string {
	return viper.GetString(logLevel)
}

func LogOutput() string {
	return viper.GetString(logOutput)
}

func GoEnvironment() string {
	return viper.GetString(goEnvironment)
}

func BatchSize() string {
	return viper.GetString(batchSize)
}

func PipelineBatchSize() int64 {
	return viper.GetInt64(pipelineBatchSize)
}

func PingInterval() string {
	return viper.GetString(pingInterval)
}

func SyncInterval() string {
	return viper.GetString(syncInterval)
}

func ProcessInterval() string {
	return viper.GetString(processInterval)
}

func CleanupInterval() string {
	return viper.GetString(cleanupInterval)
}

func FirstBlockHeight() types.Height {
	return types.Height(viper.GetInt64(firstBlockHeight))
}

func CleanupThreshold() int64 {
	return viper.GetInt64(cleanupThreshold)
}

func DbName() string {
	return viper.GetString(dbName)
}

func DbUser() string {
	return viper.GetString(dbUser)
}

func DbHost() string {
	return viper.GetString(dbHost)
}

func DbPassword() string {
	return viper.GetString(dbPassword)
}

func DbSSLMode() string {
	return viper.GetString(dbSSLMode)
}

func DbDetailedLog() bool {
	return viper.GetBool(dbDetailedLog)
}
