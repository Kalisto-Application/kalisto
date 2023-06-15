package environment

import (
	"kalisto/src/models"
	"testing"

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
	}

	err := s.e.Save(e1)
	s.NoError(err)

	envs := s.e.Get()
	s.EqualValues(map[models.EnvKind]models.Envs{"profile": {
		e1,
	}}, envs)

	// add a var and update a var to the existing env
	e1 = models.Env{
		Kind: models.EnvKindProfile,
		Name: "global",
		Vars: []models.Var{{
			Name:  "v1",
			Value: "asdf",
		}, {
			Name:  "v2",
			Value: "asdf",
		}, {
			Name:  "v3",
			Value: "asdf",
		}},
	}

	err = s.e.Save(e1)
	s.NoError(err)

	envs = s.e.Get()
	s.EqualValues(map[models.EnvKind]models.Envs{"profile": {
		e1,
	}}, envs)

	// add new env to global
	e2 := models.Env{
		Kind: models.EnvKindProfile,
		Name: "anotherName",
		Vars: []models.Var{{
			Name:  "v1",
			Value: "qwer",
		}, {
			Name:  "v2",
			Value: "qwer",
		}},
	}

	err = s.e.Save(e2)
	s.NoError(err)

	envs = s.e.Get()
	s.EqualValues(map[models.EnvKind]models.Envs{"profile": {
		e1, e2,
	}}, envs)

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
	}

	err = s.e.Save(e3)
	s.NoError(err)

	envs = s.e.Get()
	s.EqualValues(map[models.EnvKind]models.Envs{"profile": {
		e1, e2,
	}, "global": {e3}}, envs)
}

func TestEnvironment(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}
