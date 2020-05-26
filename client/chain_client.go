package client

import (
	"context"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"google.golang.org/grpc"
)

var (
	_ ChainClient = (*chainClient)(nil)
)

type ChainClient interface {
	//Queries
	GetHead() (*chainpb.GetCurrentResponse, error)
}

func NewChainClient(conn *grpc.ClientConn) *chainClient {
	return &chainClient{
		client: chainpb.NewChainServiceClient(conn),
	}
}

type chainClient struct {
	client chainpb.ChainServiceClient
}

func (r *chainClient) GetHead() (*chainpb.GetCurrentResponse, error) {
	ctx := context.Background()

	return r.client.GetCurrent(ctx, &chainpb.GetCurrentRequest{})
}