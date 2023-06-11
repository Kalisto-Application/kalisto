package compiler_test

import (
	"kalisto/src/proto/compiler"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/reflect/protoreflect"
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

	s.EqualValues(fileRegistry.NumFiles(), 1)
	fileRegistry.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		s.EqualValues(fd.Package(), "kalisto.tests.examples.service")
		s.EqualValues(fd.Services().Len(), 1)
		s.EqualValues(fd.Services().Get(0).Name(), "BookStore")
		s.EqualValues(fd.Services().Get(0).Methods().Len(), 1)
		s.EqualValues(fd.Services().Get(0).Methods().Get(0).Name(), "GetBook")

		return true
	})
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
