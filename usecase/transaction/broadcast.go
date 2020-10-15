package transaction

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
)

type broadcastUseCase struct {
	db     *store.Store
	client *client.Client
}

func NewBroadcastUseCase(db *store.Store, c *client.Client) *broadcastUseCase {
	return &broadcastUseCase{
		db:     db,
		client: c,
	}
}

func (uc *broadcastUseCase) Execute(txRaw string) (*bool, error) {
	res, err := uc.client.Transaction.Broadcast(txRaw)
	if err != nil {
		return nil, err
	}

	txSubmitted := res.GetSuccess()

	return &txSubmitted, nil
}
