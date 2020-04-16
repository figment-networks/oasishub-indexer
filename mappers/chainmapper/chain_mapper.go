package chainmapper

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasishub-indexer/models/chain"
	"github.com/figment-networks/oasishub-indexer/types"
)

func FromProxy(res chainpb.GetCurrentResponse) *chain.Model {
	return &chain.Model{
		Id:          res.Chain.Id,
		Name:        res.Chain.Name,
		GenesisTime: types.NewTimeFromTimestamp(*res.Chain.GenesisTime),
		Height:      res.Chain.Height,
	}
}
