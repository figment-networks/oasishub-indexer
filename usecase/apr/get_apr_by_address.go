package apr

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"sort"
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
	summaries, err := uc.db.BalanceSummary.GetSummariesByInterval(types.IntervalDaily, address, start, end)
	if err != nil {
		return res, err
	}

	monthlySummaries := make(map[string]MonthlyAprTotal)
	for _, summary := range summaries {
		rawAccount, err := uc.client.Account.GetByAddress(address, summary.StartHeight)
		if err != nil {
			return res, err
		}

		dailyApr, err := NewDailyApr(summary, rawAccount)
		if err != nil {
			return res, err
		}

		monthIndex := fmt.Sprintf("%d_%d", summary.TimeBucket.Year(), summary.TimeBucket.Month())
		m, ok := monthlySummaries[monthIndex]
		if ok {
			m.AprSum = m.AprSum + dailyApr.APR
			m.AprCount = m.AprCount + 1
			monthlySummaries[monthIndex] = m
		} else {
			n := MonthlyAprTotal{dailyApr.APR, 1}
			monthlySummaries[monthIndex] = n
		}
	}

	keys := make([]string, 0, len(monthlySummaries))
	for k := range monthlySummaries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	aprs := make([]MonthlyAprView, 0, len(keys))
	for i := range keys {
		a := MonthlyAprView{
			MonthInfo: keys[i],
			AvgApr:    monthlySummaries[keys[i]].AprSum / float64(monthlySummaries[keys[i]].AprCount),
		}
		aprs = append(aprs, a)
	}

	res.Result = aprs
	return res, nil
}
