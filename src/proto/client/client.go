package client

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func (c *Client) Invoke(ctx context.Context, method string, req, resp interface{}, md *metadata.MD) (string, error) {
	if err := c.conn.Invoke(ctx, method, req, resp, grpc.Header(md)); err != nil {
		st := status.Convert(err)
		if st == nil {
			return "", err
		}
		errBody, err := json.Marshal(st.Proto())
		if err != nil {
			return "", err
		}
		return string(errBody), nil
	}

	return "", nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
