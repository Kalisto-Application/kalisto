package environment

import (
	"kalisto/src/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EnvironmentTestSuite struct {
	suite.Suite

	e *Environment
}

func (s *EnvironmentTestSuite) SetupTest() {
	store := NewMockStore(s.T())
	store.EXPECT().SaveEnvs(mock.Anything).Return(nil)
	store.EXPECT().Envs().Return(nil, nil)
	e, err := NewEnvironment(store)
	s.Require().NoError(err)
	s.e = e
}

func (s *EnvironmentTestSuite) TestEnvironment() {
	w1 := uuid.NewString()
	w2 := uuid.NewString()

	// create env with 2 vars
	e1 := models.Env{
		Kind: models.EnvKindProfile,
		Name: "global",
		Vars: []models.Var{{
			Name:  "v1",
			Value: "qwer",
		}, {
			Name:  "v2",
			Value: "qwer",
		}},
		WorkspaceID: w1,
	}

	e1, err := s.e.Save(e1)
	s.NoError(err)

	envs := s.e.GetByWorkspace(w1)
	s.EqualValues(models.Envs{e1}, envs)

	// add a var and update env 1
	e1.Vars = []models.Var{{
		Name:  "v1",
		Value: "asdf",
	}, {
		Name:  "v2",
		Value: "asdf",
	}, {
		Name:  "v3",
		Value: "asdf",
	}}
	e1, err = s.e.Save(e1)
	s.NoError(err)

	envs = s.e.GetByWorkspace(w1)
	s.EqualValues(models.Envs{e1}, envs)

	// add new env 2 to workspace 1
	e2 := models.Env{
		Kind: models.EnvKindProfile,
		Name: "global",
		Vars: []models.Var{{
			Name:  "v1",
			Value: "qwer",
		}, {
			Name:  "v2",
			Value: "qwer",
		}},
		WorkspaceID: w1,
	}

	e2, err = s.e.Save(e2)
	s.NoError(err)

	envs = s.e.GetByWorkspace(w1)
	s.EqualValues(models.Envs{e2, e1}, envs)

	// add env with the same name, but another kind

	e3 := models.Env{
		Kind: models.EnvKindGlobal,
		Name: "global",
		Vars: []models.Var{{
			Name:  "v1",
			Value: "zxcv",
		}, {
			Name:  "v2",
			Value: "zxcv",
		}},
		WorkspaceID: w1,
	}
	e3, err = s.e.Save(e3)
	s.NoError(err)

	envs = s.e.GetByWorkspace(w1)
	s.EqualValues(models.Envs{e3, e2, e1}, envs)

	// create env to workspace 2

	e4 := models.Env{
		Kind: models.EnvKindGlobal,
		Name: "global",
		Vars: []models.Var{{
			Name:  "v1",
			Value: "ffff",
		}, {
			Name:  "v2",
			Value: "ffff",
		}},
		WorkspaceID: w2,
	}
	e4, err = s.e.Save(e4)
	s.NoError(err)

	envs = s.e.GetByWorkspace(w1)
	s.EqualValues(models.Envs{e3, e2, e1}, envs)
	envs = s.e.GetByWorkspace(w2)
	s.EqualValues(models.Envs{e4}, envs)
}

func TestEnvironment(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}
