package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, c Config) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.DialContext(ctx, c.Addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("proto client: failed to dial addr=%s: %w", c.Addr, err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, method string, req, resp interface{}) error {
	return c.conn.Invoke(ctx, method, req, resp)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
