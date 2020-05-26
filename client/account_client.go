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
	//Queries
	GetByPublicKey(string) (*accountpb.GetByPublicKeyResponse, error)
}

func NewAccountClient(conn *grpc.ClientConn) *accountClient {
	return &accountClient{
		client: accountpb.NewAccountServiceClient(conn),
	}
}

type accountClient struct {
	client accountpb.AccountServiceClient
}

func (r *accountClient) GetByPublicKey(key string) (*accountpb.GetByPublicKeyResponse, error) {
	ctx := context.Background()

	return r.client.GetByPublicKey(ctx, &accountpb.GetByPublicKeyRequest{PublicKey: key})
}