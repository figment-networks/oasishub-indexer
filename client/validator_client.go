package client

import (
	"context"
	"google.golang.org/grpc"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
)

var (
	_ ValidatorClient = (*validatorClient)(nil)
)

type ValidatorClient interface {
	GetByHeight(int64) (*validatorpb.GetByHeightResponse, error)
}

func NewValidatorClient(conn *grpc.ClientConn) ValidatorClient {
	return &validatorClient{
		client: validatorpb.NewValidatorServiceClient(conn),
	}
}

type validatorClient struct {
	client validatorpb.ValidatorServiceClient
}

func (r *validatorClient) GetByHeight(h int64) (*validatorpb.GetByHeightResponse, error) {
	ctx := context.Background()

	return r.client.GetByHeight(ctx, &validatorpb.GetByHeightRequest{Height: h})
}

