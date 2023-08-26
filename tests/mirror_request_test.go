package tests

import (
	"kalisto/src/assembly"
	"kalisto/src/models"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	_ "embed"
)

//go:embed testdata/request1.js
var request1 []byte

//go:embed testdata/request2.js
var request2 []byte

//go:embed testdata/request3.js
var request3 []byte

type MirrorRequestSuite struct {
	suite.Suite
}

func (s *MirrorRequestSuite) TestMirrorRequest() {
	meta := `{"content-type": 'application/grpc', "authorization": 'super token'}`
	for _, tt := range []struct {
		name string
		req  []byte
	}{
		{name: "request1", req: request1},
		{name: "request2", req: request2},
		{name: "request3", req: request3},
	} {
		s.Run(tt.name, func() {
			app, err := assembly.NewApp()
			s.Require().NoError(err)
			api := app.Api

			wd, err := os.Getwd()
			s.Require().NoError(err)
			dir := path.Join(wd, "examples/proto/")
			ws, err := api.CreateWorkspace("name", dir)
			s.Require().NoError(err)

			response, err := api.SendGrpc(models.Request{
				Addr:        "localhost:9000",
				WorkspaceID: ws.ID,
				Method:      "kalisto.tests.examples.service.BookStore.GetBook",
				Body:        string(tt.req),
				Meta:        meta,
			})
			s.Require().NoError(err)

			responseMirror, err := api.SendGrpc(models.Request{
				Addr:        "localhost:9000",
				WorkspaceID: ws.ID,
				Method:      "kalisto.tests.examples.service.BookStore.GetBook",
				Body:        string(response.Body),
				Meta:        response.MetaData,
			})
			s.Require().NoError(err)

			s.Equal(response.Body, responseMirror.Body)
			s.Equal(response.MetaData, responseMirror.MetaData)
		})
	}
}

func TestMirrorRequest(t *testing.T) {
	suite.Run(t, new(MirrorRequestSuite))
}
