package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	TaskNameBlockParser      = "BlockParser"
	TaskNameValidatorsParser = "ValidatorsParser"
	TaskNameBalanceParser    = "BalanceParser"
)

var (
	_ pipeline.Task = (*blockParserTask)(nil)
	_ pipeline.Task = (*validatorsParserTask)(nil)
)

func NewBlockParserTask() *blockParserTask {
	return &blockParserTask{
		metricObserver: indexerTaskDuration.WithLabels(TaskNameBlockParser),
	}
}

type blockParserTask struct {
	metricObserver metrics.Observer
}

type ParsedBlockData struct {
	TransactionsCount int64
	ProposerEntityUID string
}

func (t *blockParserTask) GetName() string {
	return TaskNameBlockParser
}

func (t *blockParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

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
	return &validatorsParserTask{
		metricObserver: indexerTaskDuration.WithLabels(TaskNameValidatorsParser),
	}
}

type validatorsParserTask struct {
	metricObserver metrics.Observer
}

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
	return TaskNameValidatorsParser
}

func (t *validatorsParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageParser, t.GetName(), payload.CurrentHeight))

	fetchedValidators := payload.RawValidators
	fetchedBlock := payload.RawBlock
	fetchedStakingState := payload.RawStakingState
	rewards, _ := getRewardsAndCommission(payload.RawEscrowEvents.GetAdd(), payload.CommonPoolAddress)

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

func NewBalanceParserTask() *balanceParserTask {
	return &balanceParserTask{}
}

type balanceParserTask struct{}

func (t *balanceParserTask) GetName() string {
	return TaskNameBalanceParser
}

func (t *balanceParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageParser, t.GetName(), payload.CurrentHeight))

	fetchedValidators := payload.RawValidators
	fetchedStakingState := payload.RawStakingState

	rewards, commissions := getRewardsAndCommission(payload.RawEscrowEvents.GetAdd(), payload.CommonPoolAddress)

	var err error
	balanceEvents := []model.BalanceEvent{}
	for _, fetchedValidator := range fetchedValidators {
		escrowAddr := fetchedValidator.GetAddress()

		var currActiveEscrowBalance types.Quantity
		var currTotalShares types.Quantity
		if account, ok := fetchedStakingState.GetLedger()[escrowAddr]; ok {
			currActiveEscrowBalance = types.NewQuantityFromBytes(account.GetEscrow().GetActive().GetBalance())
			currTotalShares = types.NewQuantityFromBytes(account.GetEscrow().GetActive().GetTotalShares())
		}

		// balance and shares before rewards/commission were applied
		prevActiveEscrowBalance := currActiveEscrowBalance.Clone()
		prevTotalShares := currTotalShares.Clone()

		var comShares types.Quantity
		if com, ok := commissions[escrowAddr]; ok {
			balanceEvents = append(balanceEvents, model.BalanceEvent{
				Height:        payload.CurrentHeight,
				Address:       escrowAddr,
				EscrowAddress: escrowAddr,
				Amount:        com,
				Kind:          model.Commission,
			})

			// reverse commission deposit
			err = prevActiveEscrowBalance.Sub(com)
			if err != nil {
				return err
			}

			comShares = com.Clone()
			if err = comShares.Mul(currTotalShares); err != nil {
				return err
			}
			if err = comShares.Quo(currActiveEscrowBalance); err != nil {
				return err
			}
			// reverse shares added from commission
			if err = prevTotalShares.Sub(comShares); err != nil {
				return err
			}
		}

		if reward, ok := rewards[escrowAddr]; ok {
			if err = prevActiveEscrowBalance.Sub(reward); err != nil {
				return err
			}
		}

		delegations, ok := fetchedStakingState.GetDelegations()[escrowAddr]
		var shares types.Quantity
		if ok {
			for delAddr, d := range delegations.Entries {
				shares = types.NewQuantityFromBytes(d.Shares)

				if delAddr == escrowAddr && !comShares.IsZero() {
					// reverse shares added from commission
					if err = shares.Sub(comShares); err != nil {
						return err
					}
				}
				if shares.IsZero() {
					continue
				}

				// when reward is added to escrow balance, value of shares for each delegator goes up.
				// find reward for each delegator by subtracting previous delegator balance from current delegator balance
				prevBalance := shares.Clone()
				if err = prevBalance.Mul(prevActiveEscrowBalance); err != nil {
					return err
				}
				if err = prevBalance.Quo(prevTotalShares); err != nil {
					return err
				}

				rewardsAmt := shares.Clone()
				if err = rewardsAmt.Mul(currActiveEscrowBalance); err != nil {
					return err
				}
				if err = rewardsAmt.Quo(currTotalShares); err != nil {
					return err
				}
				if err = rewardsAmt.Sub(prevBalance); err != nil {
					return err
				}

				balanceEvents = append(balanceEvents, model.BalanceEvent{
					Height:        payload.CurrentHeight,
					Address:       delAddr,
					EscrowAddress: escrowAddr,
					Amount:        rewardsAmt,
					Kind:          model.Reward,
				})
			}
		}
	}

	payload.BalanceEvents = balanceEvents
	return nil
}

func getRewardsAndCommission(rawEvents []*eventpb.AddEscrowEvent, commonPoolAddr string) (rewards map[string]types.Quantity, commissions map[string]types.Quantity) {
	rewards = make(map[string]types.Quantity)
	commissions = make(map[string]types.Quantity)

	var escrowAcnt string
	for _, rawEvent := range rawEvents {
		escrowAcnt = rawEvent.GetEscrow()
		if rawEvent.GetOwner() != commonPoolAddr {
			// rewards only come from commonpool, so skip
			continue
		}
		newAmt := types.NewQuantityFromBytes(rawEvent.GetAmount())
		existingAmt, ok := rewards[escrowAcnt]
		if ok {
			// if there's a duplicate, then the event with the higher amount is the reward (other is commission)
			if existingAmt.Cmp(newAmt) < 0 {
				rewards[escrowAcnt] = newAmt
				commissions[escrowAcnt] = existingAmt
			} else {
				commissions[escrowAcnt] = newAmt
			}
		} else {
			rewards[escrowAcnt] = newAmt
		}
	}

	return rewards, commissions
}
