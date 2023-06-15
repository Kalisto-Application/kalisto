package spec_test

import (
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/spec"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite

	c  *compiler.FileCompiler
	f  *spec.Factory
	wd string
	r  *compiler.Registry
}

func (s *FactorySuite) SetupTest() {
	s.f = spec.NewFactory()
	s.c = compiler.NewFileCompiler()

	wd, err := os.Getwd()
	s.Require().NoError(err)
	s.wd = wd

	protoPath := path.Join(s.wd, "..", "..", "..", "tests/examples/proto/service.proto")

	c := compiler.NewFileCompiler()
	fileRegistry, err := c.Compile([]string{path.Dir(protoPath)}, []string{"service.proto"})
	s.Require().NoError(err)
	s.r = fileRegistry
}

func (s *FactorySuite) TestSingleFileWithNoDeps() {
	given, err := s.f.FromRegistry(s.r)
	s.NoError(err)
	s.EqualValues(models.Spec{
		Services: []models.Service{
			{
				Name:     "BookStore",
				FullName: "kalisto.tests.examples.service.BookStore",
				Package:  "kalisto.tests.examples.service",
				Methods: []models.Method{
					{
						Name:     "GetBook",
						FullName: "kalisto.tests.examples.service.BookStore.GetBook",
						Kind:     models.CommunicationKindSimple,
						RequestMessage: models.Message{
							Name:     "GetBookRequest",
							FullName: "kalisto.tests.examples.service.GetBookRequest",
							Fields: []models.Field{
								{
									Name:         "id",
									FullName:     "kalisto.tests.examples.service.GetBookRequest.id",
									Type:         models.DataTypeString,
									DefaultValue: `""`,
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
