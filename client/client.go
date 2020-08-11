package client

import (
	"google.golang.org/grpc"
)

func New(connStr string) (*Client, error) {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,

		Account:     NewAccountClient(conn),
		Chain:       NewChainClient(conn),
		Block:       NewBlockClient(conn),
		Event:       NewEventClient(conn),
		State:       NewStateClient(conn),
		Validator:   NewValidatorClient(conn),
		Transaction: NewTransactionClient(conn),
	}, nil
}

type Client struct {
	conn *grpc.ClientConn

	Account     AccountClient
	Chain       ChainClient
	Block       BlockClient
	Event       EventClient
	State       StateClient
	Validator   ValidatorClient
	Transaction TransactionClient
}

func (c *Client) Close() error {
	return c.conn.Close()
}
