package account

import (
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
)

type getSummariesUseCase struct {
	db     *store.Store
	client *client.Client
}

func NewGetSummariesUseCase(db *store.Store, c *client.Client) *getSummariesUseCase {
	return &getSummariesUseCase{
		db:     db,
		client: c,
	}
}

type dataRow struct {
	dayStart time.Time
	account  *accountpb.Account
}

func (uc *getSummariesUseCase) Execute(address string, start, end time.Time) (DailyBalanceViewResult, error) {
	dayStart := start
	var rows []dataRow

	for {
		if dayStart == end {
			break
		}

		syncable, err := uc.db.Syncables.GetSyncableForMinTime(dayStart)
		if err != nil {
			return DailyBalanceViewResult{}, err
		}

		resp, err := uc.client.Account.GetByAddress(address, syncable.Height)
		if err != nil {
			return DailyBalanceViewResult{}, err
		}

		rows = append(rows, dataRow{
			dayStart: dayStart,
			account:  resp.GetAccount(),
		})

		dayStart = dayStart.Add(time.Hour * 24)
		if dayStart.After(end) {
			dayStart = end
		}
	}

	return toDailyBalanceViewResult(rows), nil
}
