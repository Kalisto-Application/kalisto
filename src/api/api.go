package api

import (
	"context"
	"fmt"
	"kalisto/src/environment"
	"kalisto/src/filesystem"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/interpretator"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"
	"time"

	"github.com/jhump/protoreflect/dynamic"
)

type Client interface {
	Invoke(ctx context.Context, method string, req, resp interface{}) error
	Close() error
}

type Api struct {
	compiler    *compiler.FileCompiler
	specFactory *spec.Factory
	workspace   *workspace.Workspace
	env         *environment.Environment
	newClient   func(ctx context.Context, addr string) (Client, error)
}

func New(
	compiler *compiler.FileCompiler,
	specFactory *spec.Factory,
	workspace *workspace.Workspace,
	env *environment.Environment,
	newClient func(ctx context.Context, addr string) (Client, error),
) *Api {
	return &Api{
		compiler:    compiler,
		specFactory: specFactory,
		workspace:   workspace,
		env:         env,
		newClient:   newClient,
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
	protoFiles, err := filesystem.SearchProtoFiles(request.ProtoPath)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to search proto files: %w", err)
	}

	registry, err := a.compiler.Compile([]string{protoFiles.AbsoluteDirPath}, protoFiles.RelativeProtoPaths)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	service := registry.Descriptors[0].FindService(request.FullServiceName)
	if service == nil {
		return models.Response{}, fmt.Errorf("api: failed to find service %s: %w", request.FullServiceName, err)
	}

	method := service.FindMethodByName(request.MethodName)
	if method == nil {
		return models.Response{}, fmt.Errorf("api: failed to find method %s: %w", request.MethodName, err)
	}

	c, err := a.newClient(context.TODO(), request.Addr)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create client: %w", err)
	}
	defer func() { _ = c.Close() }()

	req, err := interpretator.CreateMessageFromScript(request.Script, method.GetInputType())
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	resp := dynamic.NewMessage(method.GetOutputType())

	err = c.Invoke(ctx, "/"+service.GetFullyQualifiedName()+"/"+method.GetName(), req, resp)
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
