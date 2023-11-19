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
		kind  models.WorkspaceKind
		err   error
		errAs error
	}

	wd, err := os.Getwd()
	s.Require().NoError(err)

	for _, tt := range []testCase{
		// {
		// 	name: "proto1",
		// 	dirs: []string{path.Join(wd, "examples/proto")},
		// 	kind: models.WorkspaceKindProto,
		// },
		// {
		// 	name: "buf",
		// 	dirs: []string{path.Join(wd, "examples/buf/workspace")},
		// 	kind: models.WorkspaceKindProto,
		// },
		// {
		// 	name: "stupid links, but all the files",
		// 	dirs: []string{
		// 		path.Join(wd, "examples/buf/workspace/observabilityapi/api/v2"),
		// 		path.Join(wd, "examples/buf/workspace/observabilitytypes"),
		// 	},
		// 	kind: models.WorkspaceKindProto,
		// },
		// {
		// 	name: "no buf, just direct links",
		// 	dirs: []string{
		// 		path.Join(wd, "examples/buf/workspace/observabilityapi"),
		// 		path.Join(wd, "examples/buf/workspace/observabilitytypes"),
		// 	},
		// 	kind: models.WorkspaceKindProto,
		// },
		// {
		// 	name: "openapi",
		// 	dirs: []string{path.Join(wd, "examples/openapi/petstore.json")},
		// 	kind: models.WorkspaceKindOpenapi,
		// },
		// {
		// 	name: "openapi uber",
		// 	dirs: []string{path.Join(wd, "examples/openapi/uber.json")},
		// 	kind: models.WorkspaceKindOpenapi,
		// },
		// {
		// 	name: "openapi minimal",
		// 	dirs: []string{path.Join(wd, "examples/openapi/petstore-minimal.json")},
		// 	kind: models.WorkspaceKindOpenapi,
		// },
		{
			name: "openapi expanded",
			dirs: []string{path.Join(wd, "examples/openapi/petstore-expanded.json")},
			kind: models.WorkspaceKindOpenapi,
		},
		// {
		// 	name: "openapi simple",
		// 	dirs: []string{path.Join(wd, "examples/openapi/petstore-simple.json")},
		// 	kind: models.WorkspaceKindOpenapi,
		// },
		// {
		// 	name: "openapi with external docs",
		// 	dirs: []string{path.Join(wd, "examples/openapi/petstore-with-external-docs.json")},
		// 	kind: models.WorkspaceKindOpenapi,
		// },
		// {
		// 	name: "openapi directory",
		// 	dirs: []string{path.Join(wd, "examples/openapi/petstore/spec/swagger.json")},
		// 	kind: models.WorkspaceKindOpenapi,
		// },
	} {
		s.Run(tt.name, func() {
			app, err := assembly.NewApp(xdg.DataHome + "/kalisto/db-test-" + s.T().Name())
			s.Require().NoError(err)

			ws, err := app.Api.CreateWorkspaceV2(tt.name, tt.dirs, tt.kind)
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
