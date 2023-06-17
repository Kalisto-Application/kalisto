package workspace

import (
	"fmt"
	"kalisto/src/models"
	"sort"
	"time"

	"github.com/bobg/go-generics/v2/slices"
	"github.com/google/uuid"
)

type Store interface {
	SaveWorkspaces(w []models.Workspace) error
	Workspace() ([]models.Workspace, error)
}

type Workspace struct {
	s     Store
	cache []models.Workspace
}

func New(s Store) (*Workspace, error) {
	cache, err := s.Workspace()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace cache: %w", err)
	}
	return &Workspace{s: s, cache: cache}, nil
}

func (w *Workspace) Save(workspace models.Workspace) (models.Workspace, error) {
	workspace.ID = uuid.NewString()
	workspace.LastUsage = time.Now()
	w.cache = append(w.cache, workspace)
	return workspace, w.s.SaveWorkspaces(w.cache)
}

func (w *Workspace) Rename(id string, name string) (err error) {
	for i, workspace := range w.cache {
		if workspace.ID == id {
			workspace.Name = name
			w.cache[i] = workspace
			if err := w.s.SaveWorkspaces(w.cache); err != nil {
				return fmt.Errorf("failed to rename workspace: %w", err)
			}

			return nil
		}
	}

	return models.ErrWorkspaceNotFound
}

func (w *Workspace) Delete(id string) error {
	for i, workspace := range w.cache {
		if workspace.ID == id {
			w.cache = slices.RemoveN[[]models.Workspace](w.cache, i, 1)
			return w.s.SaveWorkspaces(w.cache)
		}
	}

	return nil
}

func (w *Workspace) Find(id string) (models.Workspace, error) {
	for i, workspace := range w.cache {
		if workspace.ID == id {
			workspace.LastUsage = time.Now()
			w.cache[i] = workspace
			if err := w.s.SaveWorkspaces(w.cache); err != nil {
				return workspace, fmt.Errorf("failed to update workspace cache: %w", err)
			}

			return workspace, nil
		}
	}

	return models.Workspace{}, models.ErrWorkspaceNotFound
}

func (w *Workspace) List() []models.Workspace {
	sort.Slice(w.cache, func(i, j int) bool {
		return w.cache[i].LastUsage.After(w.cache[j].LastUsage)
	})
	return w.cache
}
