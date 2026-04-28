package spicedb

import (
	"context"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// Replace with real gRPC call
func (c *Client) CheckPermission(ctx context.Context, subject, resource, permission string) (bool, error) {
	// TODO: integrate real SpiceDB client
	return true, nil
}
