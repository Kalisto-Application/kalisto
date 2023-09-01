package tests

import (
	"kalisto/src/assembly"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ImportSuite struct {
	suite.Suite
}

func (s *ImportSuite) TestImport() {
	type testCase struct {
		name string
		dirs []string
		err  error
	}

	wd, err := os.Getwd()
	s.Require().NoError(err)

	for _, tt := range []testCase{
		{
			name: "buf",
			dirs: []string{path.Join(wd, "examples/buf/workspace")},
		},
		{
			name: "stupid links, but all the files",
			dirs: []string{
				path.Join(wd, "examples/buf/workspace/observabilityapi/api/v2"),
				path.Join(wd, "examples/buf/workspace/observabilitytypes"),
			},
		},
		{
			name: "no buf, just direct links",
			dirs: []string{
				path.Join(wd, "examples/buf/workspace/observabilityapi"),
				path.Join(wd, "examples/buf/workspace/observabilitytypes"),
			},
		},
	} {
		s.Run(tt.name, func() {
			app, err := assembly.NewApp()
			s.Require().NoError(err)

			ws, err := app.Api.CreateWorkspace(tt.name, tt.dirs)
			if tt.err != nil {
				s.ErrorIs(err, tt.err)
			} else {
				s.Greater(len(ws.Spec.Services), 0)
			}

		})
	}
}

func TestImport(t *testing.T) {
	suite.Run(t, new(ImportSuite))
}
