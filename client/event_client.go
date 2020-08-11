package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
)

var (
	_ EventClient = (*eventClient)(nil)
)

type EventClient interface {
	GetRewardsByHeight(int64) (*eventpb.GetByHeightResponse, error)
}

func NewEventClient(conn *grpc.ClientConn) EventClient {
	return &eventClient{
		client: eventpb.NewEventServiceClient(conn),
	}
}

type eventClient struct {
	client eventpb.EventServiceClient
}

func (r *eventClient) GetRewardsByHeight(h int64) (*eventpb.GetByHeightResponse, error) {
	ctx := context.Background()

	return r.client.GetRewardsByHeight(ctx, &eventpb.GetByHeightRequest{Height: h})
}
