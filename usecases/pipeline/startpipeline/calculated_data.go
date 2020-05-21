package startpipeline

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
	"math/big"
)

var (
	_ pipeline.AsyncTask = (*CalculatedValidatorsData)(nil)
	_ pipeline.AsyncTask = (*CalculatedAccountsData)(nil)
	_ pipeline.AsyncTask = (*CalculatedEntitiesData)(nil)
)

// VALIDATORS
type CalculatedValidatorsData struct {
	items map[string]CalculatedValidatorData
}

func NewCalculatedValidatorsData() *CalculatedValidatorsData {
	return &CalculatedValidatorsData{
		items: make(map[string]CalculatedValidatorData),
	}
}

type CalculatedValidatorData struct {
	EntityUID          types.PublicKey
	NodeUID            types.PublicKey
	ConsensusUID       types.PublicKey
	Address            string
	Proposed           bool
	VotingPower        types.VotingPower
	TotalShares        types.Quantity
	PrecommitValidated int64
	PrecommitType      *string
	PrecommitIndex     *int64
}

func (data CalculatedValidatorsData) Run(errCh chan <- error, p pipeline.Payload) {
	payload := p.(*payload)
	validatorsData, err := syncablemapper.UnmarshalValidatorsData(payload.ValidatorsSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}
	blockData, err := syncablemapper.UnmarshalBlockData(payload.BlockSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}
	stateData, err := syncablemapper.UnmarshalStateData(payload.StateSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}

	for i, rv := range validatorsData {
		key := rv.Node.EntityId
		calculatedData := CalculatedValidatorData{
			EntityUID:    types.PublicKey(rv.Node.EntityId),
			NodeUID:      types.PublicKey(rv.Id),
			ConsensusUID: types.PublicKey(rv.Node.Consensus.Id),
			Address:      rv.Address,
			VotingPower:  types.VotingPower(rv.VotingPower),
		}

		// Get precommit data
		var validated int64
		var index int64
		var pType string
		if len(blockData.LastCommit.Votes) > 0 {
			// Account for situation when there is more validators than precommits
			// It means that last x validators did not have chance to vote. In that case set validated to null.
			if i > len(blockData.LastCommit.Votes)-1 {
				index = int64(i)
				validated = 2
			} else {
				precommit := blockData.LastCommit.Votes[i]

				if precommit == nil {
					validated = 0
					index = int64(i)
				} else {
					validated = 1
					index = precommit.ValidatorIndex
					pType = precommit.Type
				}
			}
		} else {
			validated = 2
		}

		calculatedData.PrecommitValidated = validated
		calculatedData.PrecommitIndex = &index
		calculatedData.PrecommitType = &pType

		// Get proposed
		calculatedData.Proposed = blockData.Header.ProposerAddress == calculatedData.Address

		// Get total shares
		delegations := stateData.Staking.Delegations[rv.Node.EntityId]
		totalShares := big.NewInt(0)
		for _, d := range delegations.Entries {
			shares := types.NewQuantityFromBytes(d.Shares)
			totalShares = totalShares.Add(totalShares, &shares.Int)
		}
		calculatedData.TotalShares = types.NewQuantity(totalShares)

		data.items[key] = calculatedData
	}
	payload.CalculatedValidatorsData = data.items
	errCh <- nil
}

// ACCOUNTS
type CalculatedAccountsData struct {
	items map[string]CalculatedAccountData
}

func NewCalculatedAccountsData() *CalculatedAccountsData {
	return &CalculatedAccountsData{
		items: make(map[string]CalculatedAccountData),
	}
}

type CalculatedAccountData struct {
	PublicKey                        types.PublicKey
	RecentGeneralBalance             types.Quantity
	RecentGeneralNonce               types.Nonce
	RecentEscrowActiveBalance        types.Quantity
	RecentEscrowActiveTotalShares    types.Quantity
	RecentEscrowDebondingBalance     types.Quantity
	RecentEscrowDebondingTotalShares types.Quantity
}

func (data CalculatedAccountsData) Run(errCh chan <- error, p pipeline.Payload) {
	payload := p.(*payload)
	stateData, err := syncablemapper.UnmarshalStateData(payload.StateSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}

	for publicKey, info := range stateData.Staking.Ledger {
		key := publicKey
		calculatedData := CalculatedAccountData{
			PublicKey:                        types.PublicKey(publicKey),
			RecentGeneralBalance:             types.NewQuantityFromBytes(info.General.Balance),
			RecentGeneralNonce:               types.Nonce(info.General.Nonce),
			RecentEscrowActiveBalance:        types.NewQuantityFromBytes(info.Escrow.Active.Balance),
			RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(info.Escrow.Active.TotalShares),
			RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(info.Escrow.Debonding.Balance),
			RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(info.Escrow.Debonding.TotalShares),
		}

		data.items[key] = calculatedData
	}
	payload.CalculatedAccountsData = data.items
	errCh <- nil
}

// ENTITIES
type CalculatedEntitiesData struct {
	items map[string]CalculatedEntityData
}

func NewCalculatedEntitiesData() *CalculatedEntitiesData {
	return &CalculatedEntitiesData{
		items: make(map[string]CalculatedEntityData),
	}
}

type CalculatedEntityData struct {
	EntityUID types.PublicKey
	isValidator bool
}

func (data CalculatedEntitiesData) Run(errCh chan <- error, p pipeline.Payload) {
	payload := p.(*payload)
	validatorsData, err := syncablemapper.UnmarshalValidatorsData(payload.ValidatorsSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}
	stateRawData, err := syncablemapper.UnmarshalStateData(payload.StateSyncable.Data)
	if err != nil {
		errCh <- err
		return
	}

	for _, entity := range stateRawData.Registry.Entities {
		key := entity.PublicKey

		isValidator := false
		for _, d := range validatorsData {
			if d.Node.EntityId == entity.PublicKey {
				isValidator = true
				break
			}
		}

		calculatedData := CalculatedEntityData{
			EntityUID: types.PublicKey(entity.PublicKey),
			isValidator: isValidator,
		}
		data.items[key] = calculatedData
	}
	payload.CalculatedEntitiesData = data.items
	errCh <- nil
}