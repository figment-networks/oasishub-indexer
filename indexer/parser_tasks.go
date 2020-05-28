package indexing

import (
	"context"
	"math/big"
	"reflect"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/types"
)

var (
	_ pipeline.Task = (*parseBlockTask)(nil)
	_ pipeline.Task = (*parseValidatorsTask)(nil)
)

func NewParseBlockTask() *parseBlockTask {
	return &parseBlockTask{}
}

type parseBlockTask struct{}

type ParsedBlockData struct {
	TransactionsCount int64
	ProposerEntityUID string
}

func (t *parseBlockTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)
	fetchedBlock := payload.RawBlock
	fetchedTransactions := payload.RawTransactions
	fetchedValidators := payload.RawValidators

	parsedBlockData := ParsedBlockData{}

	// Get transactions successCount
	parsedBlockData.TransactionsCount = int64(len(fetchedTransactions))

	// Get Proposer EntityUID
	for _, validator := range fetchedValidators {
		pa := fetchedBlock.GetHeader().GetProposerAddress()

		if pa == validator.Address {
			parsedBlockData.ProposerEntityUID = string(validator.Node.EntityId)
		}
	}
	payload.ParsedBlock = parsedBlockData
	return nil
}

func NewParseValidatorsTask() *parseValidatorsTask {
	return &parseValidatorsTask{}
}

type parseValidatorsTask struct{}

type ParsedValidatorsData map[string]parsedValidator

type parsedValidator struct {
	Proposed             bool
	PrecommitValidated   *bool
	PrecommitBlockIdFlag int64
	PrecommitIndex       int64
	TotalShares          types.Quantity
}

func (t *parseValidatorsTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)
	fetchedValidators := payload.RawValidators
	fetchedBlock := payload.RawBlock
	fetchedState := payload.RawState

	parsedData := make(ParsedValidatorsData)
	for i, rv := range fetchedValidators {
		key := rv.Node.EntityId
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
			if i > len(votes) - 1 {
				index = int64(i)
				blockIdFlag = 3
			} else {
				precommit := votes[i]
				isValidated := precommit.BlockIdFlag == 2
				validated = &isValidated
				index = precommit.ValidatorIndex
				blockIdFlag = precommit.BlockIdFlag
			}
		} else {
			index = int64(i)
			blockIdFlag = 3
		}

		calculatedData.PrecommitValidated = validated
		calculatedData.PrecommitIndex = index
		calculatedData.PrecommitBlockIdFlag = blockIdFlag

		// Get proposed
		calculatedData.Proposed = fetchedBlock.GetHeader().GetProposerAddress() == rv.Address

		// Get total shares
		delegations, ok := fetchedState.GetStaking().GetDelegations()[rv.Node.EntityId]
		totalShares := big.NewInt(0)
		if ok {
			for _, d := range delegations.Entries {
				shares := types.NewQuantityFromBytes(d.Shares)
				totalShares = totalShares.Add(totalShares, &shares.Int)
			}
		}
		calculatedData.TotalShares = types.NewQuantity(totalShares)

		parsedData[key] = calculatedData
	}
	payload.ParsedValidators = parsedData
	return nil
}
