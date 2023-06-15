package models

type Var struct {
	Name   string      `json:"name"`
	Type   DataType    `json:"type"`
	Value  interface{} `json:"value"`
	Active bool        `json:"active"`
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
	Kind   EnvKind `json:"kind"`
	Name   string  `json:"name"`
	Active bool    `json:"active"`
	Vars   []Var   `json:"vars"`
}

type Envs []Env

func (ee *Envs) Save(e Env) {
	for i := range *ee {
		if (*ee)[i].Name == e.Name {
			(*ee)[i] = e
			return
		}
	}

	*ee = append(*ee, e)
}
