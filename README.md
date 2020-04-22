### Description

This is a prototype of an app that is responsible for syncing, indexing Oasis blockchain data. 
It also exposes a JSON API for the clients to access this data.
Processing pipeline is used to sync and index data. Currently we have 3 stages in the pipeline:
* Syncer - it is responsible for getting raw data from the node
* Sequencer - it is responsible for indexing data that we want to store for every height. This data later on 
can be used to get data for specific time intervals ie. get average block-time every hour.
* Aggregator - it is responsible for indexing only current height data ie. Current balances of accounts.

### Syncables
Syncables are collections of data retrieved from Oasis node for given height. Currently we use 4 syncables:
* Block - retrieve data about block and votes
* State - retrieve data primarily about registry, staking and consensus
* Transactions - retrieve data about transactions
* Validators - retrieve data about validators and voting power

Above syncables are created during Syncer stage of the processing pipeline and stored to database.
They are also used in Sequencer and Aggregator stages to compose sequences and aggregates

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
* Account
* Entity

### Internal dependencies:
This package connects via gRPC to a oasishub-proxy which in turn connects to Oasis node.
This is required because for now the only way to connect to Oasis node is via unix socket.
oasishub-proxy also servers as a anti-corruption layer which is responsible for translating raw node 
data to the format indexer understands.

### External Packages used:
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
* GET    `/current_height`     --> get the height of the most recently synced and indexed data

### Running app

We provide a docker-compose file to make running the app easier. Below are the steps for starting the app:

1. Make sure that you have oasishub-proxy running and connected to Oasis node.
2. Run database migrations
```shell script
docker-compose run --rm migrate up
```
This will apply migrations to Timescaldb database
3. Update *PROXY_URL* in docker-compose.yml to point to your oasishub-proxy server
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

