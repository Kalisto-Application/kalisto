package api

import (
	"fmt"
	"kalisto/src/environment"
	"kalisto/src/filesystem"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"
)

type Api struct {
	compiler    *compiler.FileCompiler
	specFactory *spec.Factory
	workspace   *workspace.Workspace
	env         *environment.Environment
}

func New(
	compiler *compiler.FileCompiler,
	specFactory *spec.Factory,
	workspace *workspace.Workspace,
	env *environment.Environment,
) *Api {
	return &Api{
		compiler:    compiler,
		specFactory: specFactory,
		workspace:   workspace,
		env:         env,
	}
}

// WORKSPACE API

func (a *Api) NewWorkspace(path string) (models.Workspace, error) {
	protoFiles, err := filesystem.SearchProtoFiles(path)
	if err != nil {
		return models.Workspace{}, fmt.Errorf("api: failed to search proto files: %w", err)
	}

	registry, err := a.compiler.Compile([]string{protoFiles.AbsoluteDirPath}, protoFiles.RelativeProtoPaths)
	if err != nil {
		return models.Workspace{}, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	spc, err := a.specFactory.FromRegistry(registry)
	if err != nil {
		return models.Workspace{}, fmt.Errorf("api: failed to create spec from registry: %w", err)
	}

	ws, err := a.workspace.Save(models.Workspace{
		Name:     "New workspace",
		Spec:     spc,
		BasePath: path,
	})
	if err != nil {
		return ws, fmt.Errorf("api: failed to save workspace: %w", err)
	}

	return ws, nil
}

func (s *Api) RenameWorkspace(id string, name string) error {
	return s.workspace.Rename(id, name)
}

func (s *Api) DeleteWorkspace(id string) error {
	return s.workspace.Delete(id)
}

// ENVIRONMENT API

func (s *Api) SaveEnvironment(env models.EnvRaw) (models.Env, error) {
	vars := []models.Var{}
	return s.env.Save(models.EnvFromRaw(env, vars))
}

func (s *Api) DeleteEnvivonment(id string, workspaceID string) error {
	return s.env.Delete(id, workspaceID)
}

func (s *Api) EnvironmentsByWorkspace(id string) models.Envs {
	return s.env.GetByWorkspace(id)
}

// GRPC API

func (a *Api) SendGrpc(request models.Request) (models.Response, error) {
	return models.Response{
		Body: `{"name": "My super book"}`,
	}, nil
}
