# Description

This is a prototype of an app that is responsible for syncing data and exposing API.
It is written in Golang.

### Packages used:
* gin - Http server
* gorm - ORM with PostgreSQL interface
* cron - Cron jobs runner
* cobra - CLI builder
* viper - Configuration management
* zap - logging 

### Available Commands:

``$ oasishub migrate `` - run migrations

``$ oasishub serve `` - start http server

``$ oasishub cron start`` - start syncing process

``$ oasishub cron stop `` - stop syncing process

``$ oasishub cron inspect`` - inspect running cron jobs

``$ oasishub run calculateBlockTimes [--startTime "2020-01-22T00:00:00Z"]`` - run calculateBlockTimes task


### Suggested improvements
* Add migration files using gormigrate
* Replace cron module with module that stores queue information to persisted storage like Redis
* Split syncing service from API and have 2 independent services, one is doing only writes (syncer) to the database and other
only reads from the database (API).
* Add oasis-node and and oasis-rpc-proxy to docker-compose.yml to be able to run entire project from Docker
