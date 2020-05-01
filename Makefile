.PHONY: generate-mocks test

generate-mocks:
	@echo "[mockgen] generating mocks"
	@mockgen -destination mock/repos/accountaggrepo/account_agg_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/accountaggrepo DbRepo
	@mockgen -destination mock/repos/blockseqrepo/block_seq_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/blockseqrepo DbRepo
	@mockgen -destination mock/repos/debondingdelegationseqrepo/debonding_delegation_seq_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo DbRepo
	@mockgen -destination mock/repos/delegationseqrepo/delegation_seq_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo DbRepo
	@mockgen -destination mock/repos/validatoraggrepo/entity_agg_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/validatoraggrepo DbRepo
	@mockgen -destination mock/repos/stakingseqrepo/staking_seq_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo DbRepo
	@mockgen -destination mock/repos/syncablerepo/syncable_proxy_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/syncablerepo ProxyRepo
	@mockgen -destination mock/repos/syncablerepo/syncable_db_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/syncablerepo DbRepo
	@mockgen -destination mock/repos/transactionseqrepo/transaction_seq_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo DbRepo
	@mockgen -destination mock/repos/validatorseqrepo/validator_seq_repo_mock.go github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo DbRepo

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...