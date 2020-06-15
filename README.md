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
* Staking (disabled)
* Transactions (disabled)
* Validators
* Delegations (disabled)
* Debonding delegations (disabled)

### Aggregates
This is data that is sored in the database for the most current "entity". Aggregates are used for data
that does not change frequently or we don't care much about previous values. 
Currently we have below aggregates:
* Account (disabled)
* Validator

### Internal dependencies:
This package connects via gRPC to a oasishub-proxy which in turn connects to Oasis node.
This is required because for now the only way to connect to Oasis node is via unix socket.
oasis-rpc-proxy also servers as a anti-corruption layer which is responsible for translating raw node 
data to the format indexer understands.

### External Packages:
* `oasis-rpc-proxy` - Go proxy to Oasis node
* `indexing-engine` - A backbone for indexing process
* `gin` - Http server
* `gorm` - ORM with PostgreSQL interface
* `cron` - Cron jobs runner
* `zap` - logging 

### Environmental variables:

* `APP_ENV` - application environment (development | production) 
* `PROXY_URL` - url to oasis-rpc-proxy
* `SERVER_ADDR` - address to use for API
* `SERVER_PORT` - port to use for API
* `FIRST_BLOCK_HEIGHT` - height of first block in chain
* `INDEX_WORKER_INTERVAL` - index interval for worker
* `SUMMARIZE_WORKER_INTERVAL` - summary interval for worker
* `PURGE_WORKER_INTERVAL` - purge interval for worker
* `DEFAULT_BATCH_SIZE` - syncing batch size. Setting this value to 0 means no batch size
* `DATABASE_DSN` - PostgreSQL database URL
* `DEBUG` - turn on db debugging mode
* `LOG_LEVEL` - level of log
* `LOG_OUTPUT` - log output (ie. stdout or /tmp/logs.json)
* `ROLLBAR_ACCESS_TOKEN` - Rollbar access token for error reporting
* `ROLLBAR_SERVER_ROOT` - Rollbar server root for error reporting
* `INDEXER_METRIC_ADDR` - Prometheus server address for indexer metrics 
* `SERVER_METRIC_ADDR` - Prometheus server address for server metrics 
* `METRIC_SERVER_URL` - Url at which metrics will be accessible (for both indexer and server)
* `PURGE_BLOCK_INTERVAL` - Block sequence older than given interval will be purged
* `PURGE_BLOCK_HOURLY_SUMMARY_INTERVAL` - Block hourly summary records older than given interval will be purged
* `PURGE_BLOCK_DAILY_SUMMARY_INTERVAL` - Block daily summary records older than given interval will be purged
* `PURGE_VALIDATOR_INTERVAL` - Validator sequence older than given interval will be purged
* `PURGE_VALIDATOR_HOURLY_SUMMARY_INTERVAL` - Validator hourly summary records older than given interval will be purged
* `PURGE_VALIDATOR_DAILY_SUMMARY_INTERVAL` - Validator daily summary records older than given interval will be purged
* `INDEXER_TARGETS_FILE` - JSON file with targets and its task names 

### Available endpoints:

| Method | Path                               | Description                                                 | Params                                                                                                                                                |
|--------|------------------------------------|-------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET    | `/health`                            | health endpoint                                             | -                                                                                                                                                     |
| GET    | `/current_height`                    | get the height of the most recently synced and indexed data |                                                                                                                                                       |
| GET    | `/block`                             | return block by height                                      | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/block_times/:limit`                | get last x block times                                      | limit (required) - limit of blocks                                                                                                                    |
| GET    | `/block_summary`                     | get block summary                                           | interval (required) - time interval [hourly or daily] period (required) - summary period [ie. 24 hours]                                               |
| GET    | `/transactions`                      | get list of transactions                                    | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/validators`                        | get list of validators                                      | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/staking`                           | get staking details                                         | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/delegations`                       | get delegations                                             | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/debonding_delegations`             | get debonding delegations                                   | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/accounts`                          | get account details                                         | public_key (required) - public key of account height (optional) - height [Default: 0 = last]                                                          |
| GET    | `/validators/for_min_height/:height` | get the list of validators for height greater than provided | height (required) - height [Default: 0 = last]                                                                                                        |
| GET    | `/validators/by_entity_uid`         | get validator by entity UID                                 | entity_uid (required) - public key of entity    sequences_limit (optional) - number of sequences to include                                                                                                      |
| GET    | `/validators_summary`                | validator summary                                           | interval (required) - time interval [hourly or daily] period (required) - summary period [ie. 24 hours]  entity_uid (optional) - public key of entity |

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

### Running one-off commands

Start indexer:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer_start
```

Create summary tables for sequences:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer_summarize
```

Purge old data:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer_purge
```

### Running tests

To run tests with coverage you can use `test` Makefile target:
```shell script
make test
```

### Exporting metrics for scrapping
We use Prometheus for exposing metrics for indexer and for server.
Check environmental variables section on what variables to use to setup connection details to metrics scrapper.
We currently expose below metrics:
* `figment_indexer_height_success` (counter) - total number of successfully indexed heights
* `figment_indexer_height_error` (counter) - total number of failed indexed heights
* `figment_indexer_height_duration` (gauge) - total time required to index one height
* `figment_indexer_height_task_duration` (gauge) - total time required to process indexing task 
* `figment_indexer_use_case_duration` (gauge) - total time required to execute use case 
* `figment_database_query_duration` (gauge) - total time required to execute database query 
* `figment_server_request_duration` (gauge) - total time required to executre http request 


