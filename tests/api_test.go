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

type ApiSuite struct {
	suite.Suite
}

func toShortList(list []models.Workspace) []models.WorkspaceShort {
	res := make([]models.WorkspaceShort, len(list))
	for i := range res {
		res[i] = models.WorkspaceShort{
			ID:   list[i].ID,
			Name: list[i].Name,
		}
	}

	return res
}

func (s *ApiSuite) TestList() {
	wd, err := os.Getwd()
	s.Require().NoError(err)
	dirs := []string{
		path.Join(wd, "examples/buf/workspace/observabilityapi"),
		path.Join(wd, "examples/buf/workspace/observabilitytypes"),
	}

	app, err := assembly.NewApp(xdg.DataHome + "/kalisto.db/test-" + s.T().Name())
	s.Require().NoError(err)

	res, err := app.Api.WorkspaceList("")
	s.Require().NoError(err)
	for _, ws := range res.List {
		err = app.Api.DeleteWorkspace(ws.ID)
		s.Require().NoError(err)
	}
	res, err = app.Api.WorkspaceList("")
	s.Require().NoError(err)
	s.Require().Len(res.List, 0)

	ws1, err := app.Api.CreateWorkspace("1", dirs)
	s.Require().NoError(err)

	ws2, err := app.Api.CreateWorkspace("2", dirs)
	s.Require().NoError(err)

	ws3, err := app.Api.CreateWorkspace("3", dirs)
	s.Require().NoError(err)

	res, err = app.Api.WorkspaceList("")
	s.Require().NoError(err)
	s.Equal(res.List, toShortList([]models.Workspace{ws3, ws2, ws1}))

	res, err = app.Api.WorkspaceList(ws2.ID)
	s.Require().NoError(err)
	s.Equal(res.Main.ID, ws2.ID)
	ws2 = res.Main

	list := res.List
	s.Equal(list, toShortList([]models.Workspace{ws2, ws3, ws1}))

	res, err = app.Api.WorkspaceList(ws1.ID)
	s.Require().NoError(err)
	s.Equal(res.Main.ID, ws1.ID)
	ws1 = res.Main

	list = res.List
	s.Equal(list, toShortList([]models.Workspace{ws1, ws2, ws3}))

	app.OnShutdown(context.Background())
}

func TestApi(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
