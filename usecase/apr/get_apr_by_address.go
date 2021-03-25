package apr

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"math/big"
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

func (uc *getAprByAddressUseCase) Execute(address string, start, end *types.Time, includeDailies bool) (MonthlyAprViewResult, error) {
	var res MonthlyAprViewResult

	mostRecentSynced, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return res, err
	}
	if mostRecentSynced.Time.Before(end.Time) {
		end = types.NewTimeFromTime(mostRecentSynced.Time.Time)
	}

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

		r := NewRewardRate(summary, rawAccount)

		monthIndex := fmt.Sprintf("%d_%d", summary.TimeBucket.Year(), summary.TimeBucket.Month())
		m, ok := monthlySummaries[monthIndex]
		if ok {
			m.MonthlyRewardRate.Add(m.MonthlyRewardRate, &r.Rate)
			m.DayCount = m.DayCount + 1
			m.Dailies = append(m.Dailies, r)
			monthlySummaries[monthIndex] = m
		} else {
			mrr := new(big.Float)
			mrr.Copy(&r.Rate)
			dailies := make([]RewardRate, 0)
			dailies = append(dailies, r)
			first := MonthlyAprTotal{mrr, 1, dailies}
			monthlySummaries[monthIndex] = first
		}
	}

	keys := make([]string, 0, len(monthlySummaries))
	for k := range monthlySummaries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	aprs := make([]MonthlyAprView, 0, len(keys))
	for _, key := range keys {
		apr := new(big.Float)
		apr.SetString(monthlySummaries[key].MonthlyRewardRate.String())
		daysInYear := big.NewFloat(365)
		daysInMonth := new(big.Float)
		daysInMonth.SetInt64(monthlySummaries[key].DayCount)
		apr.Quo(apr, daysInMonth)
		apr.Mul(apr, daysInYear)
		r, _ := apr.Float64()
		a := MonthlyAprView{
			MonthInfo: key,
			AvgApr:    r,
		}
		if includeDailies {
			a.Dailies = monthlySummaries[key].Dailies
		}
		aprs = append(aprs, a)
	}

	res.Result = aprs
	return res, nil
}
