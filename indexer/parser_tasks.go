package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	BlockParserTaskName      = "BlockParser"
	ValidatorsParserTaskName = "ValidatorsParser"
)

var (
	_ pipeline.Task = (*blockParserTask)(nil)
	_ pipeline.Task = (*validatorsParserTask)(nil)
)

func NewBlockParserTask() *blockParserTask {
	return &blockParserTask{}
}

type blockParserTask struct{}

type ParsedBlockData struct {
	TransactionsCount int64
	ProposerEntityUID string
}

func (t *blockParserTask) GetName() string {
	return BlockParserTaskName
}

func (t *blockParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageParser, t.GetName(), payload.CurrentHeight))

	fetchedBlock := payload.RawBlock
	fetchedTransactions := payload.RawTransactions
	fetchedValidators := payload.RawValidators

	parsedBlockData := ParsedBlockData{}

	// Get transactions successCount
	parsedBlockData.TransactionsCount = int64(len(fetchedTransactions))

	// Get Proposer Address
	for _, validator := range fetchedValidators {
		pa := fetchedBlock.GetHeader().GetProposerAddress()

		if pa == validator.Address {
			parsedBlockData.ProposerEntityUID = string(validator.Node.EntityId)
		}
	}
	payload.ParsedBlock = parsedBlockData
	return nil
}

func NewValidatorsParserTask() *validatorsParserTask {
	return &validatorsParserTask{}
}

type validatorsParserTask struct{}

type ParsedValidatorsData map[string]parsedValidator

type parsedValidator struct {
	Proposed             bool
	PrecommitValidated   *bool
	PrecommitBlockIdFlag int64
	PrecommitIndex       int64
	TotalShares          types.Quantity
	ActiveEscrowBalance  types.Quantity
	Rewards              types.Quantity
}

func (t *validatorsParserTask) GetName() string {
	return ValidatorsParserTaskName
}

func (t *validatorsParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageParser, t.GetName(), payload.CurrentHeight))

	fetchedValidators := payload.RawValidators
	fetchedBlock := payload.RawBlock
	fetchedStakingState := payload.RawStakingState
	rewards := getRewardsFromEscrowEvents(payload.RawEscrowEvents, payload.CommonPoolAddress)

	const (
		NotValidated int64 = 1
		Validated    int64 = 2
		ValidatedNil int64 = 3
	)

	parsedData := make(ParsedValidatorsData)
	for i, fetchedValidator := range fetchedValidators {
		address := fetchedValidator.GetAddress()
		tendermintAddress := fetchedValidator.GetTendermintAddress()
		calculatedData := parsedValidator{}

		// Get precommit data
		votes := fetchedBlock.GetLastCommit().GetVotes()
		var index int64
		// 1 = Not validated
		// 2 = Validated
		// 3 = Validated nil
		var blockIdFlag int64
		var validated *bool
		if len(votes) > 0 {
			// Account for situation when there is more validators than precommits
			// It means that last x validators did not have chance to vote. In that case set validated to null.
			if i > len(votes)-1 {
				index = int64(i)
				blockIdFlag = ValidatedNil
			} else {
				precommit := votes[i]
				isValidated := precommit.BlockIdFlag == Validated
				validated = &isValidated
				index = precommit.ValidatorIndex
				blockIdFlag = precommit.BlockIdFlag
			}
		} else {
			index = int64(i)
			blockIdFlag = ValidatedNil
		}

		calculatedData.PrecommitValidated = validated
		calculatedData.PrecommitIndex = index
		calculatedData.PrecommitBlockIdFlag = blockIdFlag

		// Get proposed
		calculatedData.Proposed = fetchedBlock.GetHeader().GetProposerAddress() == tendermintAddress

		// Get total shares
		delegations, ok := fetchedStakingState.GetDelegations()[address]
		totalShares := big.NewInt(0)
		if ok {
			for _, d := range delegations.Entries {
				shares := types.NewQuantityFromBytes(d.Shares)
				totalShares = totalShares.Add(totalShares, &shares.Int)
			}
		}
		calculatedData.TotalShares = types.NewQuantity(totalShares)

		// Get active escrow
		account, ok := fetchedStakingState.GetLedger()[address]
		if ok {
			calculatedData.ActiveEscrowBalance = types.NewQuantityFromBytes(account.GetEscrow().GetActive().GetBalance())
		}

		// Get rewards
		if reward, ok := rewards[address]; ok {
			calculatedData.Rewards = reward
		}

		parsedData[address] = calculatedData
	}
	payload.ParsedValidators = parsedData
	return nil
}

func getRewardsFromEscrowEvents(rawEvents []*eventpb.AddEscrowEvent, commonPoolAddr string) map[string]types.Quantity {
	rewards := make(map[string]types.Quantity)
	for _, rawEvent := range rawEvents {
		if rawEvent.GetOwner() != commonPoolAddr {
			// rewards only come from commonpool, so skip
			continue
		}
		newAmt := types.NewQuantityFromBytes(rawEvent.GetAmount())
		existingAmt, ok := rewards[rawEvent.GetEscrow()]
		if ok {
			// if there's a duplicate, then the event with the higher amount is the reward (other is commission)
			if existingAmt.Cmp(newAmt) < 0 {
				rewards[rawEvent.GetEscrow()] = newAmt
			}

		} else {
			rewards[rawEvent.GetEscrow()] = newAmt
		}
	}
	return rewards
}
