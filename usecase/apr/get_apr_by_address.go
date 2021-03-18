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

func (uc *getAprByAddressUseCase) Execute(address string, start, end *types.Time) (MonthlyAprViewResult, error) {
	var res MonthlyAprViewResult
	summaries, err := uc.db.BalanceSummary.GetSummariesByInterval(types.IntervalMonthly, address, start, end)
	if err != nil {
		return res, err
	}

	var aprs []MonthlyAprView
	for _, summary := range summaries {
		rawAccount, err := uc.client.Account.GetByAddress(address, summary.StartHeight)
		if err != nil {
			return res, err
		}
		apr, err := NewMonthlyAprView(summary, rawAccount)
		if err != nil {
			return res, err
		}
		aprs = append(aprs, *apr)
	}

	res.Result = aprs
	return res, nil
}
