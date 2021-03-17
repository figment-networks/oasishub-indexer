package apr

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
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

func (uc *getAprByAddressUseCase) Execute(address string, start, end *types.Time) error {
	summaries, err := uc.db.BalanceSummary.GetSummariesByInterval(types.IntervalMonthly, address, start, end)
	if err != nil {
		return err
	}

	var aprs []MonthlyAprView
	for _, summary := range summaries {
		rawAccount, err := uc.client.Account.GetByAddress(address, summary.StartHeight)
		if err != nil {
			return err
		}
		// TODO: calculate apr
		apr := NewMonthlyAprView(summary, rawAccount)
		aprs = append(aprs, *apr)
	}
	//TODO: correct return

	//return MonthlyAprViewResult{Result: aprs}
	return nil
}
