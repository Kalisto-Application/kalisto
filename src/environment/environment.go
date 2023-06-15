package environment

import (
	"kalisto/src/models"
)

type Store interface {
	SaveEnvs(map[models.EnvKind]models.Envs) error
	Envs() (map[models.EnvKind]models.Envs, error)
}

type Environment struct {
	s     Store
	cache map[models.EnvKind]models.Envs
}

func NewEnvironment(s Store) (*Environment, error) {
	cache, err := s.Envs()
	if err != nil {
		return nil, err
	}
	if cache == nil {
		cache = make(map[models.EnvKind]models.Envs)
	}

	return &Environment{
		s:     s,
		cache: cache,
	}, nil
}

func (e *Environment) Save(env models.Env) error {
	envs := e.cache[env.Kind]
	envs.Save(env)
	e.cache[env.Kind] = envs
	return e.s.SaveEnvs(e.cache)
}

func (e *Environment) Get() map[models.EnvKind]models.Envs {
	return e.cache
}
