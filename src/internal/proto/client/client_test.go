package client_test

import (
	"context"
	"kalisto/src/internal/proto/client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := client.NewClient(ctx, client.Config{Addr: ":9000"})
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
