package api

import (
	"fmt"
	"kalisto/src/filesystem"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"

	"github.com/google/uuid"
)

type Api struct {
	compiler    *compiler.FileCompiler
	specFactory *spec.Factory
	workspace   *workspace.Workspace
}

func New(
	compiler *compiler.FileCompiler,
	specFactory *spec.Factory,
	workspace *workspace.Workspace,
) *Api {
	return &Api{
		compiler:    compiler,
		specFactory: specFactory,
		workspace:   workspace,
	}
}

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

func (s *Api) RenameWorkspace(id uuid.UUID, name string) error {
	return s.workspace.Rename(id, name)
}

func (s *Api) DeleteWorkspace(id uuid.UUID) error {
	return s.workspace.Delete(id)
}

func (a *Api) SendGrpc(request models.Request) (models.Response, error) {
	return models.Response{
		Body: `{"name": "My super book"}`,
	}, nil
}
