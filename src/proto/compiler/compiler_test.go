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

func (s *FileCompilerTestSuite) TestSingleFileWuthNoDeps() {
	protoPath := path.Join(s.wd, "..", "..", "..", "tests/examples/proto/service.proto")

	c := compiler.NewFileCompiler()
	fileRegistry, err := c.Compile([]string{path.Dir(protoPath)}, []string{protoPath})
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

func (s *FileCompilerTestSuite) TestMutipleFilesWithNoDeps() {

}

func (s *FileCompilerTestSuite) TestMutipleFilesWithSameLevelDep() {

}

func (s *FileCompilerTestSuite) TestMutipleFilesWithInnerLevelDep() {

}

func (s *FileCompilerTestSuite) TestMutipleFilesWithUpLevelDep() {

}

func (s *FileCompilerTestSuite) TestMutipleFilesWithExternalDep() {

}

func (s *FileCompilerTestSuite) TestMutipleFilesWithCyclicDep() {

}

func (s *FileCompilerTestSuite) TestMutipleDirectories() {

}

func (s *FileCompilerTestSuite) TestMutipleDirectoriesWithInnerLevelDep() {

}

func (s *FileCompilerTestSuite) TestMutipleDirectoriesWithUpLevelDep() {

}

func (s *FileCompilerTestSuite) TestMutipleDirectoriesWithCyclicDep() {

}

func TestFileCompiler(t *testing.T) {
	suite.Run(t, new(FileCompilerTestSuite))
}
