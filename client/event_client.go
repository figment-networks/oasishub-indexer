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
	GetAddEscrowEventsByHeight(int64) (*eventpb.GetAddEscrowEventsByHeightResponse, error)
}

func NewEventClient(conn *grpc.ClientConn) EventClient {
	return &eventClient{
		client: eventpb.NewEventServiceClient(conn),
	}
}

type eventClient struct {
	client eventpb.EventServiceClient
}

func (r *eventClient) GetAddEscrowEventsByHeight(h int64) (*eventpb.GetAddEscrowEventsByHeightResponse, error) {
	ctx := context.Background()

	return r.client.GetAddEscrowEventsByHeight(ctx, &eventpb.GetAddEscrowEventsByHeightRequest{Height: h})
}
