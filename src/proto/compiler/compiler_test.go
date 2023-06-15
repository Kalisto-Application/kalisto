package compiler_test

import (
	"kalisto/src/proto/compiler"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileCompilerTestSuite struct {
	suite.Suite

	c  *compiler.FileCompiler
	wd string
}

func (s *FileCompilerTestSuite) SetupTest() {
	s.c = compiler.NewFileCompiler()
	wd, err := os.Getwd()
	s.Require().NoError(err)
	s.wd = wd
}

func (s *FileCompilerTestSuite) TestSingleFileWithNoDeps() {
	protoPath := path.Join(s.wd, "..", "..", "..", "tests/examples/proto/service.proto")

	c := compiler.NewFileCompiler()
	fileRegistry, err := c.Compile([]string{path.Dir(protoPath)}, []string{"service.proto"})
	s.NoError(err)
	descs := fileRegistry.Descriptors

	s.EqualValues(len(descs), 1)
	fd := descs[0]
	s.EqualValues(fd.GetPackage(), "kalisto.tests.examples.service")
	s.EqualValues(len(fd.GetServices()), 1)
	s.EqualValues(fd.GetServices()[0].GetName(), "BookStore")
	s.EqualValues(len(fd.GetServices()[0].GetMethods()), 1)
	s.EqualValues(fd.GetServices()[0].GetMethods()[0].GetName(), "GetBook")
}

func (s *FileCompilerTestSuite) TestMultipleFilesWithNoDeps() {

}

func (s *FileCompilerTestSuite) TestMultipleFilesWithSameLevelDep() {

}

func (s *FileCompilerTestSuite) TestMultipleFilesWithInnerLevelDep() {

}

func (s *FileCompilerTestSuite) TestMultipleFilesWithUpLevelDep() {

}

func (s *FileCompilerTestSuite) TestMultipleFilesWithExternalDep() {

}

func (s *FileCompilerTestSuite) TestMultipleFilesWithCyclicDep() {

}

func (s *FileCompilerTestSuite) TestMultipleDirectories() {

}

func (s *FileCompilerTestSuite) TestMultipleDirectoriesWithInnerLevelDep() {

}

func (s *FileCompilerTestSuite) TestMultipleDirectoriesWithUpLevelDep() {

}

func (s *FileCompilerTestSuite) TestMultipleDirectoriesWithCyclicDep() {

}

func TestFileCompiler(t *testing.T) {
	suite.Run(t, new(FileCompilerTestSuite))
}
