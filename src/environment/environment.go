package environment

import (
	"kalisto/src/models"
	"sort"
	"time"

	"github.com/bobg/go-generics/v2/slices"
	"github.com/google/uuid"
)

type Store interface {
	SaveEnvs(map[string]models.Envs) error
	Envs() (map[string]models.Envs, error)
}

type Environment struct {
	s     Store
	cache map[string]models.Envs
}

func NewEnvironment(s Store) (*Environment, error) {
	cache, err := s.Envs()
	if err != nil {
		return nil, err
	}
	if cache == nil {
		cache = make(map[string]models.Envs)
	}

	return &Environment{
		s:     s,
		cache: cache,
	}, nil
}

func (e *Environment) Save(env models.Env) (models.Env, error) {
	if env.ID == "" {
		env.ID = uuid.NewString()
		env.CreatedAt = time.Now()
	}
	envs := e.cache[env.WorkspaceID]
	envs.Save(env)
	e.cache[env.WorkspaceID] = envs
	return env, e.s.SaveEnvs(e.cache)
}

func (e *Environment) GetByWorkspace(id string) models.Envs {
	envs := e.cache[id]
	sort.Slice(envs, func(i, j int) bool {
		return envs[i].CreatedAt.After(envs[j].CreatedAt)
	})
	return envs
}

func (s *Environment) Delete(id string, workspaceID string) error {
	envs := s.cache[workspaceID]
	for i, env := range envs {
		if env.ID == id {
			envs = slices.RemoveN[models.Envs](envs, i, 1)
			s.cache[workspaceID] = envs
			s.s.SaveEnvs(s.cache)
			return nil
		}
	}

	return nil
}
