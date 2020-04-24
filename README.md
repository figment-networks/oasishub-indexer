# oasishub-indexer

This is a prototype of an app that is responsible for syncing, indexing Oasis blockchain data. 
It also exposes a JSON API for the clients to access this data.

Processing pipeline is used to sync and index data. Currently we have 3 stages in the pipeline:

- **Syncer** - Responsible for getting raw data from the node
- **Sequencer** - Responsible for indexing data that we want to store for every height. This data later on 
can be used to get data for specific time intervals ie. get average block-time every hour.
- **Aggregator** - Responsible for indexing only current height data ie. Current balances of accounts.

### Syncables

Syncables are collections of data retrieved from Oasis node for given height. Currently we use 4 syncables:

- **Block**        - Retrieves data about block and votes
- **State**        - Retrieves data primarily about registry, staking and consensus
- **Transactions** - Retrieves data about transactions
- **Validators**   - Retrieves data about validators and voting power

Above syncables are created during Syncer stage of the processing pipeline and stored to database.
They are also used in Sequencer and Aggregator stages to compose sequences and aggregates

### Sequences

Sequence is a data that is stored in the database for every height. Sequences are used for data that 
changes frequently and we want to know about those changes. This data is perfect 
for displaying change over time using graphs on the front-end.

Supported sequences:

- Block
- Staking
- Transactions
- Validators
- Delegations
- Debonding delegations

### Aggregates

Aggregates is a data that is sored in the database for the most current "entity". 
Aggregates are used for data that does not change frequently or we don't care 
much about previous values.

Supported aggregates:

- Account
- Entity

### Dependencies

This package connects via gRPC to a `oasishub-proxy` which in turn connects to Oasis node.
This is required because for now the only way to connect to Oasis node is via unix socket.
oasishub-proxy also servers as a anti-corruption layer which is responsible for 
translating raw node data to the format indexer understands.

#### External Packages

- gin   - HTTP server
- gorm  - ORM with PostgreSQL interface
- cron  - Cron jobs runner
- cobra - CLI builder
- viper - Configuration management
- zap   - Logging

### Available Commands

- `$ cli`    - Executes CLI commands
- `$ server` - Starts HTTP server
- `$ job`    - Starts CRON job

### Available Endpoints

| Method | Path                            | Description
|--------|---------------------------------|------------------------------------
| GET    | /ping                           | Healthcheck endpoint, returns `pong`
| GET    | /blocks/:height                 | Returns block by height
| GET    | /block_times/:limit             | Returns last x block times
| GET    | /block_times_interval/:interval | Returns block times for specific interval ie. '5 minutes' or '1 hour'
| GET    | /transactions/:height           | Returns list of transactions for given height
| GET    | /validators/:height             | Returns list of validators for given height
| GET    | /staking/:height                | Returns staking information for given height
| GET    | /delegations/:height            | Returns delegations for given height
| GET    | /debonding_delegations/:height  | Returns debonding delegations for given height
| GET    | /accounts/:public_key           | Returns account information by public key
| GET    | /current_height                 | Returns the height of the most recently synced and indexed data

### Running app

We provide a docker-compose file to make running the app easier.
Below are the steps for starting the app:

1. Make sure that you have `oasishub-proxy` running and connected to Oasis node.
2. Run database migrations
```shell script
docker-compose run --rm migrate up
```
This will apply migrations to Timescaldb database
3. Update `PROXY_URL` in docker-compose.yml to point to your oasishub-proxy server
4. To run a job which will start indexing process:
```shell script
docker-compose up job
``` 
5. To start a web server (port 8081 is set by default):
```shell script
docker-compose up
```

### Running tests

To run tests with coverage you can use `test` Makefile target:

```shell script
make test
```