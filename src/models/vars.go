package models

import "time"

type Var struct {
	Name  string      `json:"name"`
	Type  DataType    `json:"type"`
	Value interface{} `json:"value"`
}

type EnvKind string

const (
	EnvKindGlobal    EnvKind = "global"
	EnvKindWorkspace EnvKind = "workspace"
	EnvKindProfile   EnvKind = "profile"
)

var AllEnvKinds = []EnvKind{
	EnvKindGlobal,
	EnvKindWorkspace,
	EnvKindProfile,
}

type EnvRaw struct {
	ID          string  `json:"id"`
	Kind        EnvKind `json:"kind"`
	Name        string  `json:"name"`
	Active      bool    `json:"active"`
	Vars        string  `json:"vars"`
	WorkspaceID string  `json:"workspaceID"`
}

type Env struct {
	ID          string    `json:"id"`
	Kind        EnvKind   `json:"kind"`
	Name        string    `json:"name"`
	Active      bool      `json:"active"`
	Vars        []Var     `json:"vars"`
	WorkspaceID string    `json:"workspaceID"`
	CreatedAt   time.Time `json:"createdAt"`
}

func EnvFromRaw(e EnvRaw, vars []Var) Env {
	return Env{
		ID:          e.ID,
		Kind:        e.Kind,
		Name:        e.Name,
		Active:      e.Active,
		Vars:        vars,
		WorkspaceID: e.WorkspaceID,
	}
}

type Envs []Env

func (ee *Envs) Save(e Env) {
	for i := range *ee {
		if (*ee)[i].ID == e.ID {
			(*ee)[i] = e
			return
		}
	}

	*ee = append(*ee, e)
}
