package api_test

import (
	"context"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"kalisto/src/api"
	"kalisto/src/models"
	"os"
	"path"
	"testing"

	. "github.com/onsi/gomega"
)

func TestApi_SendGrpc(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	req := models.Request{
		ProtoPath:       path.Join(wd, "..", "..", "tests/examples/proto/service.proto"),
		FullServiceName: "kalisto.tests.examples.service.BookStore",
		MethodName:      "GetBook",
		Script: `
			a = "1"
 			request = {id: a}
			`,
	}

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
	})

	// Call SendGrpc
	got, err := a.SendGrpc(req)

	// Assert
	require.NoError(t, err)

	g := NewGomegaWithT(t)
	g.Expect(got.Body).To(MatchJSON(`{"name": "All quiet on the western front"}`))

	mockClient.AssertExpectations(t)
}
