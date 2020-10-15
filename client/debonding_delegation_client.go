package client

import (
	"context"
	"google.golang.org/grpc"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/debondingdelegation/debondingdelegationpb"
)

var (
	_ DebondingDelegationClient = (*debondingDelegationClient)(nil)
)

type DebondingDelegationClient interface {
	GetByAddress(string, int64) (*debondingdelegationpb.GetByAddressResponse, error)
}

func NewDebondingDelegationClient(conn *grpc.ClientConn) DebondingDelegationClient {
	return &debondingDelegationClient{
		client: debondingdelegationpb.NewDebondingDelegationServiceClient(conn),
	}
}

type debondingDelegationClient struct {
	client debondingdelegationpb.DebondingDelegationServiceClient
}

func (r *debondingDelegationClient) GetByAddress(address string, h int64) (*debondingdelegationpb.GetByAddressResponse, error) {
	ctx := context.Background()

	return r.client.GetByAddress(ctx, &debondingdelegationpb.GetByAddressRequest{Address: address, Height: h})
}

