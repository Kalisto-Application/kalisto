package spec_test

import (
	"kalisto/src/internal/proto/compiler"
	"kalisto/src/internal/proto/spec"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type FactorySuite struct {
	suite.Suite

	c  *compiler.FileCompiler
	f  *spec.Factory
	wd string
	r  *protoregistry.Files
}

func (s *FactorySuite) SetupTest() {
	s.f = spec.NewFactory()
	s.c = compiler.NewFileCompiler()

	wd, err := os.Getwd()
	s.Require().NoError(err)
	s.wd = wd

	protoPath := path.Join(s.wd, "..", "..", "..", "..", "tests/examples/proto/service.proto")

	c := compiler.NewFileCompiler()
	fileRegistry, err := c.Compile([]string{path.Dir(protoPath)}, []string{protoPath})
	s.Require().NoError(err)
	s.r = fileRegistry
}

func (s *FactorySuite) TestSingleFileWuthNoDeps() {
	given, err := s.f.FromRegistry(s.r)
	s.NoError(err)
	s.EqualValues(spec.Spec{
		Services: []spec.Service{
			{
				Name: "BookStore",
				Methods: []spec.Method{
					{
						Name: "GetBook",
						Kind: spec.CommunicationKindSimple,
						RequestMessage: spec.Message{
							Name: "GetBookRequest",
							Fields: []spec.Field{
								{
									Name: "id",
									Type: spec.DataTypeString,
								},
							},
						},
					},
				},
			},
		},
	}, given)
}

func TestFactory(t *testing.T) {
	suite.Run(t, new(FactorySuite))
}
