package chain

import "github.com/figment-networks/oasishub-indexer/types"

type Model struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	GenesisTime *types.Time `json:"genesis_time"`
	Height      int64      `json:"height"`
}
