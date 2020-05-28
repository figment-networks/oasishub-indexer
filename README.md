### Description

This is a prototype of an app that is responsible for syncing, indexing Oasis blockchain data. 
It also exposes a JSON API for the clients to access this data.
Processing pipeline is used to sync and index data. Currently we have 3 stages in the pipeline:
* Syncer - it is responsible for getting raw data from the node
* Sequencer - it is responsible for indexing data that we want to store for every height. This data later on 
can be used to get data for specific time intervals ie. get average block-time every hour.
* Aggregator - it is responsible for indexing only current height data ie. Current balances of accounts.

### Syncables
Syncables are collections of data retrieved from Oasis node for given height.
Above syncables are created during Syncer stage of the processing pipeline and stored to database.

### Sequences
This is data that is stored in the database for every height. Sequences are used for data that 
changes frequently and we want to know about those changes. This data is perfect for displaying change over time using graphs on the front-end.
Currently we store below sequences: 
* Block
* Staking
* Transactions
* Validators
* Delegations
* Debonding delegations

### Aggregates
This is data that is sored in the database for the most current "entity". Aggregates are used for data
that does not change frequently or we don't care much about previous values. 
Currently we have below aggregates:
* Account (pending)
* Validator

### Internal dependencies:
This package connects via gRPC to a oasishub-proxy which in turn connects to Oasis node.
This is required because for now the only way to connect to Oasis node is via unix socket.
oasis-rpc-proxy also servers as a anti-corruption layer which is responsible for translating raw node 
data to the format indexer understands.

### External Packages used:
* gin - Http server
* gorm - ORM with PostgreSQL interface
* cron - Cron jobs runner
* zap - logging 

### Environmental variables:

* `APP_ENV` - application environment (development | production) 
* `PROXY_URL` - url to oasis-rpc-proxy
* `SERVER_ADDR` - address to use for API
* `SERVER_PORT` - port to use for API
* `FIRST_BLOCK_HEIGHT` - height of first block in chain
* `SYNC_INTERVAL` - data sync interval
* `DEFAULT_BATCH_SIZE` - syncing batch size. Setting this value to 0 means no batch size
* `DATABASE_DSN` - postgreSQL database URL
* `DEBUG` - turn on db debugging mode
* `LOG_LEVEL` - level of log
* `LOG_OUTPUT` - log output (ie. stdout or /tmp/logs.json)
* `ROLLBAR_ACCESS_TOKEN` - Rollbar access token for error reporting
* `ROLLBAR_SERVER_ROOT` - Rollbar server root for error reporting
* `METRIC_SERVER_ADDR` - Prometheus server address
* `METRIC_SERVER_URL` - Url at which metrics will be accessible 

### Available endpoints:

* GET    `/health`                     --> ping endpoint
* GET    `/blocks`           --> return block by height. You can pass optional height query param.
* GET    `/block_times`       --> get last x block times
* GET    `/block_times_interval` --> get block times for specific interval ie. '5 minutes' or '1 hour'
* GET    `/transactions`     --> get list of transactions for given height. You can pass optional height query param.
* GET    `/validators`       --> get list of validators for given height. You can pass optional height query param.
* GET    `/staking`          --> get staking information for given height. You can pass optional height query param.
* GET    `/delegations`      --> get delegations for given height. You can pass optional height query param.
* GET    `/debonding_delegations` --> get debonding delegations for given height. You can pass optional height query param.
* GET    `/accounts?public_key=:public_key`     --> get account information by public key
* GET    `/current_height`     --> get the height of the most recently synced and indexed data
* GET    `/validators/for_min_height/:height`     --> get the list of validators for height greater than provided
* GET    `/validators/by_entitiy_uid?entity_uid=:entity_uid`     --> get validator by entity UID
* GET    `/validators/shares_interval?entity_uid=:entity_uid&interval=:interval&period=:period`     --> get shares for validator for specific interval ie. '5 minutes' or '1 hour'
* GET    `/validators/voting_power_interval?entity_uid=:entity_uid&interval=:interval&period=:period`     --> get voting power for validator for specific interval ie. '5 minutes' or '1 hour'
* GET    `/validators/uptime_interval?entity_uid=:entity_uid&interval=:interval&period=:period`     --> get uptime for validator for specific interval ie. '5 minutes' or '1 hour'

* GET    `/validators/total_shares_interval?interval=:interval&period=:period`     --> get total shares for specific interval ie. '5 minutes' or '1 hour'
* GET    `/validators/total_voting_power_interval?interval=:interval&period=:period`     --> get voting power for specific interval ie. '5 minutes' or '1 hour'

### Running app

Once you have created a database and specified all configuration options, you
need to migrate the database. You can do that by running the command below:

```bash
oasishub-indexer -config path/to/config.json -cmd=migrate
```

Start the data indexer:

```bash
oasishub-indexer -config path/to/config.json -cmd=worker
```

Start the API server:

```bash
oasishub-indexer -config path/to/config.json -cmd=server
```

IMPORTANT!!! Make sure that you have oasishub-proxy running and connected to Oasis node.

### Running tests

To run tests with coverage you can use `test` Makefile target:
```shell script
make test
```

### Exporting metrics for scrapping
We use Prometheus for exposing metrics.
You can use `METRIC_SERVER_ADDR` and `METRIC_SERVER_URL` to setup connection details to metrics scrapper (see Environmental variables section above).
We currently expose 4 metrics:
* `figment_indexer_height_success` (counter) - total number of successfully indexed heights
* `figment_indexer_height_error` (counter) - total number of failed indexed heights
* `figment_indexer_height_duration` (gauge) - total time required to index one height
* `figment_indexer_height_task_duration` (gauge) - total time required to process indexing task 


