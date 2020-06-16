package cli

import "fmt"

var (
	appName    = "oasishub-indexer"
	appVersion = "0.4.0"
	gitCommit  = "-"
	goVersion  = "1.13"
)

func versionString() string {
	return fmt.Sprintf(
		"%s %s (git: %s, %s)",
		appName,
		appVersion,
		gitCommit,
		goVersion,
	)
}
