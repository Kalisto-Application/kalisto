package tests

import (
	"context"
	"kalisto/src/assembly"
	"kalisto/src/models"
	"os"
	"path"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/suite"
)

type ImportSuite struct {
	suite.Suite
}

func (s *ImportSuite) TestImport() {
	type testCase struct {
		name  string
		dirs  []string
		err   error
		errAs error
	}

	wd, err := os.Getwd()
	s.Require().NoError(err)

	for _, tt := range []testCase{
		// {
		// 	name: "proto1",
		// 	dirs: []string{path.Join(wd, "examples/proto")},
		// },
		// {
		// 	name: "buf",
		// 	dirs: []string{path.Join(wd, "examples/buf/workspace")},
		// },
		// {
		// 	name: "stupid links, but all the files",
		// 	dirs: []string{
		// 		path.Join(wd, "examples/buf/workspace/observabilityapi/api/v2"),
		// 		path.Join(wd, "examples/buf/workspace/observabilitytypes"),
		// 	},
		// },
		// {
		// 	name: "no buf, just direct links",
		// 	dirs: []string{
		// 		path.Join(wd, "examples/buf/workspace/observabilityapi"),
		// 		path.Join(wd, "examples/buf/workspace/observabilitytypes"),
		// 	},
		// },
		{
			name: "openapi",
			dirs: []string{path.Join(wd, "examples/openapi/swagger.yaml")},
		},
		// {
		// 	name:  "openapi dir",
		// 	dirs:  []string{path.Join(wd, "examples/openapi")},
		// 	errAs: &models.ErrorOpenapiFileCantBeDir{File: "examples/openapi"},
		// },
	} {
		s.Run(tt.name, func() {
			app, err := assembly.NewApp(xdg.DataHome + "/kalisto/db-test-" + s.T().Name())
			s.Require().NoError(err)

			ws, err := app.Api.CreateWorkspaceV2(tt.name, tt.dirs, models.WorkspaceKindOpenapi)
			if tt.err != nil {
				s.ErrorIs(err, tt.err)
			} else if tt.errAs != nil {
				s.Contains(err.Error(), tt.errAs.Error())
			} else {
				s.Greater(len(ws.Spec.Services), 0)
			}
			app.OnShutdown(context.Background())
		})
	}
}

func TestImport(t *testing.T) {
	suite.Run(t, new(ImportSuite))
}
