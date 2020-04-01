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

### Available endpoints:

* GET    `/ping`                     --> ping endpoint
* GET    `/blocks/:height`           --> return block by height
* GET    `/block_times/:limit`       --> get last x block times
* GET    `/block_times_interval/:interval` --> get block times for specific interval ie. '5 minutes' or '1 hour'
* GET    `/transactions/:height`     --> get list of transactions for given height
* GET    `/validators/:height`       --> get list of validators for given height
* GET    `/staking/:height`          --> get staking information for given height
* GET    `/delegations/:height`      --> get delegations for given height
* GET    `/debonding_delegations/:height` --> get debonding delegations for given height
* GET    `/accounts/:public_key`     --> get account information by public key
