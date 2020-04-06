package errors

import (
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/rollbar/rollbar-go"
)

func RecoverError() {
	err := recover()
	rollbar.LogPanic(err, true)
}

func init() {
	rollbar.SetToken(config.RollbarAccessToken())
	rollbar.SetEnvironment(config.GoEnvironment())
	rollbar.SetServerRoot("github.com/figment-networks/oasishub-indexer")
}
