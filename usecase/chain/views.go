package chain

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DetailsView struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	GoVersion  string `json:"go_version"`

	ChainAppVersion   uint64 `json:"chain_app_version"`
	ChainBlockVersion uint64 `json:"chain_block_version"`
	ChainID           string `json:"chain_id"`
	ChainName         string `json:"chain_name"`

	GenesisHeight int64      `json:"genesis_height"`
	GenesisTime   types.Time `json:"genesis_time"`

	LastIndexVersion  int64      `json:"last_index_version"`
	LastIndexedHeight int64      `json:"last_indexed_height"`
	LastIndexedTime   types.Time `json:"last_indexed_time"`
	LastIndexedAt     types.Time `json:"last_indexed_at"`
	Lag               int64      `json:"indexing_lag"`
}

func ToDetailsView(recentSyncable *model.Syncable, headResponse *chainpb.GetHeadResponse, statusResponse *chainpb.GetStatusResponse) *DetailsView {
	return &DetailsView{
		AppName:    config.AppName,
		AppVersion: config.AppVersion,
		GoVersion:  config.GoVersion,

		ChainID:           statusResponse.GetId(),
		ChainName:         statusResponse.GetName(),
		ChainAppVersion:   recentSyncable.AppVersion,
		ChainBlockVersion: recentSyncable.BlockVersion,

		GenesisHeight: statusResponse.GetGenesisHeight(),
		GenesisTime:   *types.NewTimeFromTimestamp(*statusResponse.GetGenesisTime()),

		LastIndexVersion:  recentSyncable.IndexVersion,
		LastIndexedHeight: recentSyncable.Height,
		LastIndexedTime:   recentSyncable.Time,
		LastIndexedAt:     recentSyncable.CreatedAt,
		Lag:               headResponse.Height - recentSyncable.Height,
	}
}
