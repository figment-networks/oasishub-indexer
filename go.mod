module github.com/figment-networks/oasishub-indexer

go 1.14

require (
	github.com/figment-networks/indexing-engine v0.1.9
	github.com/figment-networks/oasis-rpc-proxy v0.3.11
	github.com/gin-gonic/gin v1.5.0
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/common v0.10.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/rollbar/rollbar-go v1.2.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.24.0 // indirect
)
