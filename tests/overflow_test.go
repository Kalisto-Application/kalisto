package tests

import (
	"kalisto/src/assembly"
	"kalisto/src/models"
	server "kalisto/tests/examples/server_seq"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	_ "embed"
)

//go:embed testdata/overflow/float_overflow.js
var floatOverflowRequest []byte

//go:embed testdata/overflow/int32_overflow.js
var int32OverflowRequest []byte

//go:embed testdata/overflow/uint32_overflow.js
var uint32OverflowRequest []byte

//go:embed testdata/overflow/sint32_overflow.js
var sint32OverflowRequest []byte

//go:embed testdata/overflow/fixed32_overflow.js
var fixed32OverflowRequest []byte

//go:embed testdata/overflow/sfixed32_overflow.js
var sfixed32OverflowRequest []byte

//go:embed testdata/overflow/float_valid.js
var floatValidRequest []byte

//go:embed testdata/overflow/int32_valid.js
var int32ValidRequest []byte

//go:embed testdata/overflow/uint32_valid.js
var uint32ValidRequest []byte

//go:embed testdata/overflow/sint32_valid.js
var sint32ValidRequest []byte

//go:embed testdata/overflow/fixed32_valid.js
var fixed32ValidRequest []byte

//go:embed testdata/overflow/sfixed32_valid.js
var sfixed32ValidRequest []byte

type OverflowSuite struct {
	suite.Suite

	close func() error
}

func (s *OverflowSuite) SetupSuite() {
	close, _, err := server.Run(":9000")
	s.Require().NoError(err)
	s.close = close
	time.Sleep(time.Millisecond * 200)
}

func (s *OverflowSuite) TearDownSuite() {
	s.close()
}

func (s *OverflowSuite) TestOverflow() {
	meta := `{"content-type": 'application/grpc', authorization: 'super token'}`
	for _, tt := range []struct {
		name string
		req  []byte
		err  string
	}{
		{name: "valid float", req: floatValidRequest},
		{name: "overflow float", req: floatOverflowRequest, err: "float overflow"},
		{name: "valid int32", req: int32ValidRequest},
		{name: "overflow int32", req: int32OverflowRequest, err: "integer overflow"},
		{name: "valid uint32", req: uint32ValidRequest},
		{name: "overflow uint32", req: uint32OverflowRequest, err: "integer overflow"},

		{name: "valid fixed32", req: fixed32ValidRequest},
		{name: "overflow fixed32", req: fixed32OverflowRequest, err: "integer overflow"},
		{name: "valid sfixed32", req: sfixed32ValidRequest},
		{name: "overflow sfixed32", req: sfixed32OverflowRequest, err: "integer overflow"},
		{name: "valid sint32", req: sint32ValidRequest},
		{name: "overflow sint32", req: sint32OverflowRequest, err: "integer overflow"},
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

			_, err = api.SendGrpc(models.Request{
				Addr:        "localhost:9000",
				WorkspaceID: ws.ID,
				Method:      "kalisto.tests.examples.service.BookStore.GetBook",
				Body:        string(tt.req),
				Meta:        meta,
			})
			if tt.err != "" {
				s.Require().Error(err)
				if !s.True(strings.HasSuffix(err.Error(), tt.err)) {
					log.Default().Printf("given: %s\n", err.Error())
				}
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestOverflowSuite(t *testing.T) {
	suite.Run(t, new(OverflowSuite))
}
