package client

import (
	"context"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"google.golang.org/grpc"
)

var (
	_ AccountClient = (*accountClient)(nil)
)

type AccountClient interface {
	GetByAddress(string, int64) (*accountpb.GetByAddressResponse, error)
}

func NewAccountClient(conn *grpc.ClientConn) *accountClient {
	return &accountClient{
		client: accountpb.NewAccountServiceClient(conn),
	}
}

type accountClient struct {
	client accountpb.AccountServiceClient
}

func (r *accountClient) GetByAddress(address string, height int64) (*accountpb.GetByAddressResponse, error) {
	ctx := context.Background()

	return r.client.GetByAddress(ctx, &accountpb.GetByAddressRequest{Address: address, Height: height})
}