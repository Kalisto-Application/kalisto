package api

import (
	"context"
	"encoding/json"
	"fmt"
	"kalisto/src/environment"
	"kalisto/src/filesystem"
	"kalisto/src/models"
	"kalisto/src/proto/client"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/interpreter"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"
	"reflect"
	"strings"
	"time"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/grpc/metadata"
)

type Api struct {
	ctx context.Context

	compiler      *compiler.FileCompiler
	specFactory   *spec.Factory
	workspace     *workspace.Workspace
	globalVars    *environment.GlovalVars
	newClient     func(ctx context.Context, addr string) (*client.Client, error)
	protoRegistry *compiler.Descritors
}

func New(
	compiler *compiler.FileCompiler,
	specFactory *spec.Factory,
	workspace *workspace.Workspace,
	globalVars *environment.GlovalVars,
	newClient func(ctx context.Context, addr string) (*client.Client, error),
	protoRegistry *compiler.Descritors,
) *Api {
	return &Api{
		compiler:      compiler,
		specFactory:   specFactory,
		workspace:     workspace,
		globalVars:    globalVars,
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
	if len(spc.Services) == 0 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    "error",
			Title:   "Can't create a workspace",
			Message: "No services found",
		})
		return models.Workspace{}, fmt.Errorf("no services found")
	}

	ws, err := a.workspace.Save(models.Workspace{
		Name:      "New workspace",
		Spec:      spc,
		BasePath:  path,
		TargetUrl: "localhost:9000",
		LastUsage: time.Now(),
	})
	if err != nil {
		return ws, fmt.Errorf("api: failed to save workspace: %w", err)
	}

	a.protoRegistry.Add(ws.ID, registry)

	return ws, nil
}

func (s *Api) DeleteWorkspace(id string) error {
	return s.workspace.Delete(id)
}

func (s *Api) FindWorkspaces() ([]models.Workspace, error) {
	list := s.workspace.List()
	for i, w := range list {
		w, err := s.enrichWorkspace(w, w.LastUsage)
		if err != nil {
			return nil, err
		}
		list[i] = w
	}
	return list, nil
}

func (s *Api) GetWorkspace(id string) (models.Workspace, error) {
	ws, err := s.workspace.Find(id)
	if err != nil {
		return ws, err
	}

	return s.enrichWorkspace(ws, time.Now())
}

func (s *Api) UpdateWorkspace(ws models.Workspace) error {
	return s.workspace.Update(ws)
}

func (s *Api) enrichWorkspace(ws models.Workspace, lastUsage time.Time) (models.Workspace, error) {
	registry, err := s.protoRegistryFromPath(ws.BasePath)
	if err != nil {
		// TODO: MARK AS INVALID
		return ws, err
	}
	s.protoRegistry.Add(ws.ID, registry)

	spec, err := s.specFactory.FromRegistry(registry)
	if err != nil {
		return ws, err
	}

	newWs := ws
	newWs.Spec = spec
	newWs.LastUsage = lastUsage

	if !reflect.DeepEqual(ws, newWs) {
		if err := s.workspace.Update(newWs); err != nil {
			return newWs, err
		}
	}

	return newWs, nil
}

// ENVIRONMENT API

func (s *Api) GetGlobalVars() string {
	return s.globalVars.Get()
}

func (s *Api) SaveGlovalVars(vars string) error {
	if err := s.globalVars.Save(vars); err != nil {
		return err
	}
	ip := interpreter.NewInterpreter("")
	if _, err := ip.Raw(vars); err != nil {
		return err
	}

	return nil
}

// SCRIPTING

func (a *Api) SaveScript(WorkspaceID, script string) error {
	ws, err := a.workspace.Find(WorkspaceID)
	if err != nil {
		return err
	}
	ws.Script = script
	if err := a.workspace.Update(ws); err != nil {
		return nil
	}

	return nil
}

// GRPC API

func (a *Api) SendGrpc(request models.Request) (models.Response, error) {
	if strings.TrimSpace(request.Body) == "" {
		return models.Response{}, nil
	}

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

	vars := a.GetGlobalVars()
	ip := interpreter.NewInterpreter(vars)
	specInputMessage, err := ws.Spec.FindInputMessage(sd.GetFullyQualifiedName(), md.GetName())
	if err != nil {
		return models.Response{}, err
	}

	req, err := ip.CreateMessageFromScript(request.Body, md.GetInputType(), ws.Spec, specInputMessage)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create request: %w", err)
	}

	meta, err := ip.CreateMetadata(request.Meta)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create metadata: %w", err)
	}

	resp := dynamic.NewMessage(md.GetOutputType())

	ctx = metadata.NewOutgoingContext(ctx, meta)
	responseMeta := metadata.MD{}
	apiErr, err := c.Invoke(ctx, "/"+sd.GetFullyQualifiedName()+"/"+md.GetName(), req, resp, &responseMeta)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to invoke method: %w", err)
	}

	specOutputMessage, err := ws.Spec.FindOutputMessage(sd.GetFullyQualifiedName(), md.GetName())
	if err != nil {
		return models.Response{}, err
	}
	var body string
	if apiErr != "" {
		body = apiErr
	} else {
		body, err = a.specFactory.MessageAsJsString(specOutputMessage, resp)
		if err != nil {
			return models.Response{}, fmt.Errorf("api: failed to present response as js object: %w", err)
		}
	}
	metaJson, err := json.Marshal(responseMeta)
	if err != nil {
		return models.Response{}, err
	}

	return models.Response{
		Body:     body,
		MetaData: string(metaJson),
	}, nil
}

func (a *Api) RunScript(request models.ScriptCall) (string, error) {
	if strings.TrimSpace(request.Body) == "" {
		return "", nil
	}

	reg, err := a.protoRegistry.Get(request.WorkspaceID)
	if err != nil {
		return "", err
	}

	ws, err := a.workspace.Find(request.WorkspaceID)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	c, err := a.newClient(ctx, request.Addr)
	if err != nil {
		return "", fmt.Errorf("api: failed to create client: %w", err)
	}
	defer func() { _ = c.Close() }()

	vars := a.GetGlobalVars()
	ip := interpreter.NewInterpreter(vars)

	resp, err := ip.RunScript(ctx, request.Body, ws.Spec, reg, c)
	if err != nil {
		return "", fmt.Errorf("api: failed to create request: %w", err)
	}
	b, err := resp.MarshalJSONIndent()
	if err != nil {
		return "", fmt.Errorf("api: failed to marshal response: %w", err)
	}
	return string(b), nil
}

func (s *Api) protoRegistryFromPath(path string) (*compiler.Registry, error) {
	protoFiles, err := filesystem.SearchProtoFiles(path)
	if err != nil {
		return nil, fmt.Errorf("api: failed to search proto files: %w", err)
	}

	registry, err := s.compiler.Compile(protoFiles.AbsoluteDirPath, protoFiles.RelativeProtoPaths, protoFiles.BufDirs)
	if err != nil {
		return nil, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	return registry, nil
}
