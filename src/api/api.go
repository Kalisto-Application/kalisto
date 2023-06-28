package api

import (
	"context"
	"fmt"
	"kalisto/src/environment"
	"kalisto/src/filesystem"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/interpreter"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"
	"time"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Client interface {
	Invoke(ctx context.Context, method string, req, resp interface{}) error
	Close() error
}

type Api struct {
	ctx context.Context

	compiler      *compiler.FileCompiler
	specFactory   *spec.Factory
	workspace     *workspace.Workspace
	env           *environment.Environment
	newClient     func(ctx context.Context, addr string) (Client, error)
	protoRegistry *compiler.Descritors
}

func New(
	compiler *compiler.FileCompiler,
	specFactory *spec.Factory,
	workspace *workspace.Workspace,
	env *environment.Environment,
	newClient func(ctx context.Context, addr string) (Client, error),
	protoRegistry *compiler.Descritors,
) *Api {
	return &Api{
		compiler:      compiler,
		specFactory:   specFactory,
		workspace:     workspace,
		env:           env,
		newClient:     newClient,
		protoRegistry: protoRegistry,
	}
}

func SetContext(a *Api, ctx context.Context) {
	a.ctx = ctx
}

// WORKSPACE API

func (a *Api) NewWorkspace() (models.Workspace, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return models.Workspace{}, err
	}

	registry, err := a.protoRegistryFromPath(path)
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

	a.protoRegistry.Add(ws.ID, registry)

	return ws, nil
}

func (s *Api) RenameWorkspace(id string, name string) error {
	return s.workspace.Rename(id, name)
}

func (s *Api) DeleteWorkspace(id string) error {
	return s.workspace.Delete(id)
}

func (s *Api) FindWorkspaces() ([]models.Workspace, error) {
	list := s.workspace.List()
	for _, w := range list {
		registry, err := s.protoRegistryFromPath(w.BasePath)
		if err != nil {
			// TODO: MARK AS INVALID
			continue
		}
		s.protoRegistry.Add(w.ID, registry)
	}
	return list, nil
}

func (s *Api) GetWorkspace(id string) (models.Workspace, error) {
	ws, err := s.workspace.Find(id)
	if err != nil {
		return ws, err
	}

	registry, err := s.protoRegistryFromPath(ws.BasePath)
	if err != nil {
		// TODO: MARK AS INVALID
		return ws, err
	}
	s.protoRegistry.Add(ws.ID, registry)

	return ws, nil
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
	reg, err := a.protoRegistry.Get(request.WorkspaceID)
	if err != nil {
		return models.Response{}, err
	}
	sd, md, err := reg.FindMethod(models.MethodName(request.Method))
	if err != nil {
		return models.Response{}, err
	}

	ws, err := a.workspace.Find(request.WorkspaceID)
	if err != nil {
		return models.Response{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	c, err := a.newClient(ctx, request.Addr)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create client: %w", err)
	}
	defer func() { _ = c.Close() }()

	req, err := interpreter.CreateMessageFromScript(request.Body, md.GetInputType(), ws.Spec, sd.GetFullyQualifiedName(), md.GetName())
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create request: %w", err)
	}

	md.GetService()
	resp := dynamic.NewMessage(md.GetOutputType())

	err = c.Invoke(ctx, "/"+sd.GetFullyQualifiedName()+"/"+md.GetName(), req, resp)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to invoke method: %w", err)
	}

	b, err := resp.MarshalJSONIndent()
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to marshal response: %w", err)
	}

	return models.Response{
		Body: string(b),
	}, nil
}

func (s *Api) protoRegistryFromPath(path string) (*compiler.Registry, error) {
	protoFiles, err := filesystem.SearchProtoFiles(path)
	if err != nil {
		return nil, fmt.Errorf("api: failed to search proto files: %w", err)
	}

	registry, err := s.compiler.Compile([]string{protoFiles.AbsoluteDirPath}, protoFiles.RelativeProtoPaths)
	if err != nil {
		return nil, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	return registry, nil
}
