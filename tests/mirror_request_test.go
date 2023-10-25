package tests

import (
	"context"
	"kalisto/src/assembly"
	"kalisto/src/models"
	"kalisto/tests/examples/server"
	"os"
	"path"
	"testing"
	"time"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/suite"

	_ "embed"
)

//go:embed testdata/request1.js
var request1 []byte

//go:embed testdata/request2.js
var request2 []byte

//go:embed testdata/request3.js
var request3 []byte

//go:embed testdata/request4.js
var request4 []byte

type MirrorRequestSuite struct {
	suite.Suite

	close func() error
}

func (s *MirrorRequestSuite) SetupSuite() {
	close, _, err := server.Run(":9000")
	s.Require().NoError(err)
	s.close = close
	time.Sleep(time.Millisecond * 200)
}

func (s *MirrorRequestSuite) TearDownSuite() {
	s.close()
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
		{name: "request4", req: request4},
	} {
		s.Run(tt.name, func() {
			app, err := assembly.NewApp(xdg.DataHome + "/kalisto.db/test-" + s.T().Name())
			s.Require().NoError(err)
			api := app.Api

			wd, err := os.Getwd()
			s.Require().NoError(err)
			dir := path.Join(wd, "examples/proto/")
			ws, err := api.CreateWorkspace("name", []string{dir})
			s.Require().NoError(err)

			response, err := api.SendGrpc(models.Request{
				Addr:        "localhost:9000",
				WorkspaceID: ws.ID,
				Method:      "kalisto.tests.examples.service.BookStore.Mirror",
				Body:        string(tt.req),
				Meta:        meta,
			})
			s.Require().NoError(err)

			responseMirror, err := api.SendGrpc(models.Request{
				Addr:        "localhost:9000",
				WorkspaceID: ws.ID,
				Method:      "kalisto.tests.examples.service.BookStore.Mirror",
				Body:        string(response.Body),
				Meta:        response.MetaData,
			})
			s.Require().NoError(err)

			AssertJsObjectsAreEqual(s.T(), response.Body, responseMirror.Body)
			s.EqualValues(response.MetaData, responseMirror.MetaData)

			app.OnShutdown(context.Background())
		})
	}
}

func TestMirrorRequest(t *testing.T) {
	suite.Run(t, new(MirrorRequestSuite))
}
