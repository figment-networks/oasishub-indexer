package apr

import (
	"fmt"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
)

var (
	timeFormat = "2006-01-02"
)

type getAprByAddressUseCase struct {
	db     *store.Store
	client *client.Client
}

func NewGetAprByAddressUseCase(db *store.Store, c *client.Client) *getAprByAddressUseCase {
	return &getAprByAddressUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getAprByAddressUseCase) Execute(address string, start, end *types.Time) ([]DailyApr, error) {
	mostRecentSynced, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return []DailyApr{}, err
	}
	if mostRecentSynced.Time.Before(end.Time) {
		end = types.NewTimeFromTime(mostRecentSynced.Time.Time)
	}

	rewardSeqs, err := uc.db.BalanceSummary.GetSummariesByInterval(types.IntervalDaily, address, start, end)
	if err != nil {
		return []DailyApr{}, err
	}

	if len(rewardSeqs) == 0 {
		return []DailyApr{}, fmt.Errorf("No rewards exist for account: %w", http.ErrNotFound)
	}

	rewardLookup := make(map[string]model.BalanceSummary)
	delegationLookup := make(map[string]*delegationpb.GetByAddressResponse)
	validators := []string{}

	for _, r := range rewardSeqs {
		rewardLookupKey := fmt.Sprintf("%s.%s", r.EscrowAddress, r.TimeBucket.Format(timeFormat))
		rewardLookup[rewardLookupKey] = r

		delegationLookupKey := fmt.Sprintf("%s.%d", r.Address, r.StartHeight)
		if _, ok := delegationLookup[delegationLookupKey]; !ok {
			delegation, err := uc.client.Delegation.GetByAddress(r.Address, r.StartHeight) // fetches shares
			if err != nil {
				return []DailyApr{}, err
			}
			delegationLookup[delegationLookupKey] = delegation
		}

		validators = append(validators, r.EscrowAddress)
	}

	summaries, err := uc.db.ValidatorSummary.FindAllByTimePeriod(start, end, validators...)
	if err != nil {
		return []DailyApr{}, err
	}

	return toAPRView(summaries, rewardLookup, delegationLookup)
}
