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

type Env struct {
	ID          string    `json:"id"`
	Kind        EnvKind   `json:"kind"`
	Name        string    `json:"name"`
	Vars        string    `json:"vars"`
	WorkspaceID string    `json:"workspaceID"`
	CreatedAt   time.Time `json:"createdAt"`
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
