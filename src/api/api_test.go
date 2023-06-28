package api_test

import (
	"context"
	"kalisto/src/api"
	"kalisto/src/models"
	"testing"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	. "github.com/onsi/gomega"
)

func TestApi_SendGrpc(t *testing.T) {
	req := models.Request{}

	// Create a new mock client
	mockClient := new(api.MockClient)
	mockClient.On("Invoke", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			// this argument is a pointer to the response
			resPtr := args.Get(3)
			if res, ok := resPtr.(*dynamic.Message); ok {
				res.SetFieldByName("name", "All quiet on the western front")
			} else {
				t.Error("could not cast response to dynamic.Message")
			}
		}).Return(nil)

	mockClient.On("Close").Return(nil)

	// Initialize API with mock client
	a := api.New(nil, nil, nil, nil, func(ctx context.Context, addr string) (api.Client, error) {
		return mockClient, nil
	}, nil)

	// Call SendGrpc
	got, err := a.SendGrpc(req)

	// Assert
	require.NoError(t, err)

	g := NewGomegaWithT(t)
	g.Expect(got.Body).To(MatchJSON(`{"name": "All quiet on the western front"}`))

	mockClient.AssertExpectations(t)
}
