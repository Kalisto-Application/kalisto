package tests

import (
	"kalisto/src/assembly"
	"kalisto/src/models"
	server "kalisto/tests/examples/server_seq"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	_ "embed"
)

//go:embed testdata/script.js
var script []byte

type SequenceScriptSuite struct {
	suite.Suite

	close func() error
}

func (s *SequenceScriptSuite) SetupSuite() {
	close, _, err := server.Run(":9000")
	s.Require().NoError(err)
	s.close = close
	time.Sleep(time.Millisecond * 200)
}

func (s *SequenceScriptSuite) TearDownSuite() {
	s.close()
}

func (s *SequenceScriptSuite) TestSequenceScript() {
	meta := `{"content-type": 'application/grpc', authorization: 'super token'}`
	for _, tt := range []struct {
		name string
		req  []byte
	}{
		{name: "sequence script", req: script},
	} {
		s.Run(tt.name, func() {
			app, err := assembly.NewApp()
			s.Require().NoError(err)
			api := app.Api

			wd, err := os.Getwd()
			s.Require().NoError(err)
			dir := path.Join(wd, "examples/proto_sequence/")
			ws, err := api.CreateWorkspace("name", dir)
			s.Require().NoError(err)

			response, err := api.RunScript(models.ScriptCall{
				Addr:        "localhost:9000",
				WorkspaceID: ws.ID,
				Body:        string(tt.req),
				Meta:        meta,
			})
			s.Require().NoError(err)
			s.Require().JSONEq(response, `{"value": 3, "rpc": "Third"}`)
		})
	}
}

func TestSequenceScript(t *testing.T) {
	suite.Run(t, new(SequenceScriptSuite))
}
