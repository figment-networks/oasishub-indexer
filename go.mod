module github.com/figment-networks/oasishub-indexer

go 1.13

replace github.com/tendermint/tendermint => github.com/oasislabs/tendermint v0.32.8-oasis1

require (
	github.com/figment-networks/oasis-rpc-proxy v0.0.0-20200330190657-a4ba84fc2c07
	github.com/gin-gonic/gin v1.5.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/jinzhu/gorm v1.9.12
	github.com/lib/pq v1.3.0
	github.com/oasislabs/oasis-core/go v0.0.0-20200304114707-807935769a93
	github.com/robfig/cron/v3 v3.0.1
	github.com/rollbar/rollbar-go v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	go.uber.org/zap v1.14.0
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
)
