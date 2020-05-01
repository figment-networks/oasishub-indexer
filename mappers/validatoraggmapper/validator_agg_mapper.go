package validatoraggmapper

import (
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/validatoragg"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecases/pipeline/startpipeline"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func FromCalculatedData(height types.Height, time types.Time, validatorData startpipeline.CalculatedValidatorData, existing *validatoragg.Model) (*validatoragg.Model, errors.ApplicationError) {
	validatorAgg := validatoragg.Model{
		Aggregate: &shared.Aggregate{
			RecentAtHeight: height,
			RecentAt:       time,
		},

		EntityUID:               validatorData.EntityUID,
		RecentAddress:           validatorData.Address,
		RecentTotalShares:       validatorData.TotalShares,
		RecentVotingPower:       validatorData.VotingPower,
		RecentAsValidatorHeight: height,
	}

	if validatorData.Proposed {
		validatorAgg.RecentProposedHeight = height
	}

	if existing == nil {
		// Create
		validatorAgg.Aggregate.StartedAtHeight = height
		validatorAgg.Aggregate.StartedAt = time

		if validatorData.PrecommitValidated == 0 {
			validatorAgg.AccumulatedUptime = 0
			validatorAgg.AccumulatedUptimeCount = 1
		} else if validatorData.PrecommitValidated == 1 {
			validatorAgg.AccumulatedUptime = 1
			validatorAgg.AccumulatedUptimeCount = 1
		} else {
			// We don't count out of range as offline
			validatorAgg.AccumulatedUptime = 0
			validatorAgg.AccumulatedUptimeCount = 0
		}
	} else {
		// Update
		if validatorData.PrecommitValidated == 0 {
			validatorAgg.AccumulatedUptime = existing.AccumulatedUptime + 0
			validatorAgg.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
		} else if validatorData.PrecommitValidated == 1 {
			validatorAgg.AccumulatedUptime = existing.AccumulatedUptime + 1
			validatorAgg.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
		} else {
			// We don't count out of range as offline
			validatorAgg.AccumulatedUptime = existing.AccumulatedUptime + 0
			validatorAgg.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 0
		}
	}

	if !validatorAgg.Valid() {
		return nil, errors.NewErrorFromMessage("validator aggregate not valid", errors.NotValid)
	}
	return &validatorAgg, nil
}

type DetailsView struct {
	*shared.Model
	*shared.Aggregate

	EntityUID                types.PublicKey   `json:"entity_uid"`
	RecentAddress            string            `json:"recent_address"`
	RecentVotingPower        types.VotingPower `json:"recent_voting_power"`
	RecentTotalShares        types.Quantity    `json:"recent_total_shares"`
	RecentAsValidatorHeight  types.Height      `json:"recent_as_validator_height"`
	RecentProposedHeight     types.Height      `json:"recent_proposed_height"`
	AccumulatedProposedCount int64             `json:"accumulated_proposed_count"`
	Uptime                   float64           `json:"uptime"`

	RecentDelegations          []delegationseq.Model          `json:"recent_delegations"`
	RecentDebondingDelegations []debondingdelegationseq.Model `json:"recent_debonding_delegations"`
}

func ToDetailsView(m validatoragg.Model, currDs []delegationseq.Model, recDds []debondingdelegationseq.Model) *DetailsView {
	return &DetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		EntityUID:                m.EntityUID,
		RecentAddress:            m.RecentAddress,
		RecentVotingPower:        m.RecentVotingPower,
		RecentTotalShares:        m.RecentTotalShares,
		RecentAsValidatorHeight:  m.RecentAsValidatorHeight,
		RecentProposedHeight:     m.RecentProposedHeight,
		AccumulatedProposedCount: m.AccumulatedProposedCount,
		Uptime:                   float64(m.AccumulatedUptime) / float64(m.AccumulatedUptimeCount),

		RecentDelegations:          currDs,
		RecentDebondingDelegations: recDds,
	}
}

type ListView struct {
	Items []validatoragg.Model `json:"items"`
}

func ToListView(ms []validatoragg.Model) *ListView {
	return &ListView{
		Items: ms,
	}
}
