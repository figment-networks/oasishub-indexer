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

``$ cli `` - execute CLI commands

``$ server `` - start http server

``$ job`` - start CRON job
