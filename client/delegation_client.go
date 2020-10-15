package client

import (
	"context"
	"google.golang.org/grpc"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
)

var (
	_ DelegationClient = (*delegationClient)(nil)
)

type DelegationClient interface {
	GetByAddress(string, int64) (*delegationpb.GetByAddressResponse, error)
}

func NewDelegationClient(conn *grpc.ClientConn) DelegationClient {
	return &delegationClient{
		client: delegationpb.NewDelegationServiceClient(conn),
	}
}

type delegationClient struct {
	client delegationpb.DelegationServiceClient
}

func (r *delegationClient) GetByAddress(address string, h int64) (*delegationpb.GetByAddressResponse, error) {
	ctx := context.Background()

	return r.client.GetByAddress(ctx, &delegationpb.GetByAddressRequest{Address: address, Height: h})
}

