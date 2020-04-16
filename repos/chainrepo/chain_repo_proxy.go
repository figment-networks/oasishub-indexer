package chainrepo

import (
	"context"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasishub-indexer/mappers/chainmapper"
	"github.com/figment-networks/oasishub-indexer/models/chain"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

type ProxyRepo interface {
	//Queries
	GetCurrent() (*chain.Model, errors.ApplicationError)
}

type proxyRepo struct {
	chainClient chainpb.ChainServiceClient
}

func NewProxyRepo(chainClient chainpb.ChainServiceClient) ProxyRepo {
	return &proxyRepo{
		chainClient: chainClient,
	}
}

func (r *proxyRepo) GetCurrent() (*chain.Model, errors.ApplicationError) {
	ctx := context.Background()

	res, err := r.chainClient.GetCurrent(ctx, &chainpb.GetCurrentRequest{})
	if err != nil {
		return nil, errors.NewError("could not get chain info from proxy", errors.ProxyRequestError, err)
	}
	return chainmapper.FromProxy(*res), nil
}
