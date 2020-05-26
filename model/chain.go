package model

import "github.com/figment-networks/oasishub-indexer/types"

type Chain struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	GenesisTime *types.Time `json:"genesis_time"`
	Height      int64      `json:"height"`
}
