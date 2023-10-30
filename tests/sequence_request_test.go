package tests

import (
	"context"
	"kalisto/src/assembly"
	"kalisto/src/definitions/proto/interpreter"
	"kalisto/src/models"
	"kalisto/tests/examples/server"
	"kalisto/tests/examples/server_seq"
	"os"
	"path"
	"testing"
	"time"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/suite"

	_ "embed"
)

//go:embed testdata/script.js
var script []byte

//go:embed testdata/scriptMirror.js
var scriptMirror []byte

type SequenceScriptSuite struct {
	suite.Suite

	close func() error
}

func (s *SequenceScriptSuite) SetupSuite() {
}

func (s *SequenceScriptSuite) TestScriptSequence() {
	meta := `{"content-type": 'application/grpc', authorization: 'super token'}`
	close, closed, err := server_seq.Run(":9000")
	s.Require().NoError(err)
	time.Sleep(time.Millisecond * 200)

	app, err := assembly.NewApp(xdg.DataHome + "/kalisto.db/test-" + s.T().Name())
	s.Require().NoError(err)
	api := app.Api

	wd, err := os.Getwd()
	s.Require().NoError(err)
	dir := path.Join(wd, "examples/proto_sequence/")
	ws, err := api.CreateWorkspace("name", []string{dir})
	s.Require().NoError(err)

	response, err := api.RunScript(models.ScriptCall{
		Addr:        "localhost:9000",
		WorkspaceID: ws.ID,
		Body:        string(script),
		Meta:        meta,
	})
	s.Require().NoError(err)

	AssertJsObjectsAreEqual(s.T(), `{"value": 3, "rpc": "Third"}`, response)

	close()
	app.OnShutdown(context.Background())
	<-closed
}

func (s *SequenceScriptSuite) TestMirrorScripting() {
	meta := `{"content-type": 'application/grpc', authorization: 'super token'}`
	close, closed, err := server.Run(":9000")
	s.Require().NoError(err)
	time.Sleep(time.Millisecond * 200)

	app, err := assembly.NewApp(xdg.DataHome + "/kalisto.db/test-" + s.T().Name())
	s.Require().NoError(err)
	api := app.Api

	wd, err := os.Getwd()
	s.Require().NoError(err)
	dir := path.Join(wd, "examples/proto")
	ws, err := api.CreateWorkspace("name", []string{dir})
	s.Require().NoError(err)

	response, err := api.RunScript(models.ScriptCall{
		Addr:        "localhost:9000",
		WorkspaceID: ws.ID,
		Body:        string(scriptMirror),
		Meta:        meta,
	})
	s.Require().NoError(err)

	ip := interpreter.NewInterpreter("")
	ex1, err := ip.ExportValue(string(request1), "")
	s.Require().NoError(err)
	ex2, err := ip.ExportValue(response, "")
	s.Require().NoError(err)

	ex1["time"] = ex1["time"].(time.Time).UTC().Format(time.RFC3339)
	for k, v := range ex2 {
		if v == nil {
			delete(ex2, k)
		}
	}
	s.EqualValues(ex1, ex2)

	close()
	app.OnShutdown(context.Background())
	<-closed
}

func (s *SequenceScriptSuite) TestSequentialRequests() {
	close, closed, err := server.Run(":9000")
	s.Require().NoError(err)
	time.Sleep(time.Millisecond * 200)

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
		Method:      "kalisto.tests.examples.service.BookStore.CreateBook",
		Body:        `{name: "Some book"}`,
		Meta:        ``,
	})
	s.Require().NoError(err)

	AssertJsObjectsAreEqual(s.T(), response.Body, `{
		id: '1',
	  }`)

	response, err = api.SendGrpc(models.Request{
		Addr:        "localhost:9000",
		WorkspaceID: ws.ID,
		Method:      "kalisto.tests.examples.service.BookStore.GetBook",
		Body:        string(response.Body),
		Meta:        response.MetaData,
	})
	s.Require().NoError(err)
	AssertJsObjectsAreEqual(s.T(), response.Body, `{
		id: '1',
		name: 'Clean Code',
	  }`)

	close()
	app.OnShutdown(context.Background())
	<-closed
}

func TestSequenceScript(t *testing.T) {
	suite.Run(t, new(SequenceScriptSuite))
}
