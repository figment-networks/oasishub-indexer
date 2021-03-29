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
* `PURGE_SEQUENCES_INTERVAL` - Sequence older than given interval will be purged _[DEFAULT: 24h]_
* `PURGE_SYSTE_EVENTS_INTERVAL` - System events older than given interval will be purged _[DEFAULT: 24h]_
* `PURGE_HOURLY_SUMMARY_INTERVAL` - Hourly summaries records older than given interval will be purged _[DEFAULT: 24h]_
* `INDEXER_CONFIG_FILE` - JSON file with indexer configuration 

### Available endpoints:

| Method | Path                               | Description                                                 | Params                                                                                                                                                |
|--------|------------------------------------|-------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET    | `/health`                            | health endpoint                                             | -                                                                                                                                                     |
| GET    | `/status`                            | status of the application and chain                         | -                                                                                                                                                     |
| GET    | `/block`                             | return block by height                                      | `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/block_times/:limit`                | get last x block times                                      | `limit (required)` - limit of blocks                                                                                                                    |
| GET    | `/blocks_summary`                    | get block summary                                           | `interval (required)` - time interval [hourly or daily] `period (required)` - summary period [ie. 24 hours]                                               |
| GET    | `/transactions`                      | get list of transactions                                    | `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/staking`                           | get staking details                                         | `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/delegations`                       | get delegations                                             | `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/delegations/:address`              | get delegations for address                                 | `address (required)` - address of account    `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/debonding_delegations`             | get debonding delegations                                   | `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/debonding_delegations/:address`    | get debonding delegations for address                       | `address (required)` - address of account    `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/account/:address`                  | get account details                                         | `address (required)` - address of account `height (optional)` - height [Default: 0 = last]                                                          |
| GET    | `/validators`                        | get list of validators                                      | `height (optional)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/validators/for_min_height/:height` | get the list of validators for height greater than provided | `height (required)` - height [Default: 0 = last]                                                                                                        |
| GET    | `/validator/:address`                | get validator by address                                    | `address (required)` - validator's address    `sequences_limit (optional)` - number of sequences to include                                                                                                      |
| GET    | `/validators_summary`                | validator summary                                           | `interval (required)` - time interval [hourly or daily] `period (required)` - summary period [ie. 24 hours]  `address (optional)` - address of entity |
| GET    | `/system_events/:address`            | system events for given actor                               | `address (required)` - address of account `after (optional)` - return events after with height greater than provided height  `kind (optional)` - system event kind |
| POST   | `/transactions`                      | broadcast transaction                                       | `tx_raw (required)` - raw transaction data as string                                                                                                        |
| GET   | `/apr/:address`                      | get time series of annualized rewards rates calculated per month   | `start (required)` - start date in format `2006-01-02` `end (required)` - end date in format `2006-01-02` `address (required)` - address of account  |

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

Start indexing process:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer:index
```

Start backfill process:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer:backfill
```

Create summary tables for sequences:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer:summarize
```

Purge old data:
```bash
oasishub-indexer -config path/to/config.json -cmd=indexer:purge
```

Decorate validator aggregates:
```bash
oasishub-indexer -config path/to/config.json -cmd=validators:decorate -file=/file/to/csv
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
* `figment_server_request_duration` (gauge) - total time required to execute http request 


### Using indexer configuration file
Indexing process is configured using JSON file. A typical indexer config file looks similar to:
```json
{
  "versions": [
    {
      "id": 1,
      "parallel": false,
      "targets": [1, 2]
    },
    {
      "id": 2,
      "parallel": false,
      "targets": [3]
    }
  ],
  "shared_tasks": [
    "HeightMetaRetriever",
    "MainSyncer"
  ],
  "available_targets": [
    {
      "id": 1,
      "name": "index_block_sequences",
      "desc": "Creates and persists block sequences",
      "tasks": [
        "BlockFetcher",
        "ValidatorFetcher",
        "TransactionFetcher",
        "BlockParser",
        "BlockSeqCreator",
        "SyncerPersistor",
        "BlockSeqPersistor"
      ]
    },
    {
      "id": 2,
      "name": "index_validator_sequences",
      "desc": "Creates and persists validator sequences",
      "tasks": [
        "BlockFetcher",
        "StakingStateFetcher",
        "ValidatorFetcher",
        "ValidatorsParser",
        "ValidatorSeqCreator",
        "SyncerPersistor",
        "ValidatorSeqPersistor"
      ]
    },
    {
      "id": 3,
      "name": "index_validator_aggregates",
      "desc": "Creates and persists validator aggregates",
      "tasks": [
        "BlockFetcher",
        "StakingStateFetcher",
        "ValidatorFetcher",
        "ValidatorsParser",
        "ValidatorAggCreator",
        "SyncerPersistor",
        "ValidatorAggPersistor"
      ]
    }
  ]
}
```

In this file we have 3 main sections:
* **versions** - describes available index versions and what targets to run for every version. `parallel` option specifies if the new index version can be run in parallel.
* **shared_tasks** - tasks that are shared between targets. These tasks will be run for every target in `targets` section.
* **targets** - array of available targets along with their tasks. Target represents some specific outcome of indexing process (ie. index validators) and it tells the system what tasks have to be run to satisfy this outcome.

### Updating indexer version flow

- Add task to a stage in the `indexer` package
- If necessary, create `migration` file(s)
- Create a new target in the targets section of the `indexer_config.json`
- Add a new version in `versions` section. Inside of it, specify targets that need to be run to satisfy this version. If you provide `parallel=true` option it will mean that backfill and indexing can be run in parallel.
