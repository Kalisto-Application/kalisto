package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"kalisto/src/definitions"
	ocompiler "kalisto/src/definitions/openapi/compiler"
	"kalisto/src/definitions/proto/client"
	pcompiler "kalisto/src/definitions/proto/compiler"
	"kalisto/src/definitions/proto/interpreter"
	"kalisto/src/definitions/proto/render"
	"kalisto/src/filesystem"
	"kalisto/src/models"
	"kalisto/src/pkg/runtime"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"kalisto/src/db"

	"github.com/bufbuild/protocompile/reporter"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/dynamic"
	rpkg "github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/grpc/metadata"
)

type Api struct {
	ctx context.Context
	mx  sync.RWMutex

	protoCompiler   *pcompiler.FileCompiler
	openapiCompiler *ocompiler.Compiler

	store         *db.DB
	newClient     func(ctx context.Context, addr string) (*client.Client, error)
	registryStore *definitions.RegistryStore

	runtime runtime.Runtime
}

func New(
	protoCompiler *pcompiler.FileCompiler,
	openapiCompiler *ocompiler.Compiler,
	store *db.DB,
	newClient func(ctx context.Context, addr string) (*client.Client, error),
	registryStore *definitions.RegistryStore,
	runtime runtime.Runtime,
) *Api {
	return &Api{
		protoCompiler:   protoCompiler,
		openapiCompiler: openapiCompiler,
		store:           store,
		newClient:       newClient,
		registryStore:   registryStore,
		runtime:         runtime,
	}
}

func SetContext(a *Api, ctx context.Context) {
	a.mx.Lock()
	a.ctx = ctx
	a.mx.Unlock()
}

func (a *Api) Context() context.Context {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return a.ctx
}

func (a *Api) Runtime() runtime.Runtime {
	return a.runtime
}

// WORKSPACE API

func (a *Api) FindProtoFiles() (models.ProtoDir, error) {
	path, err := a.runtime.OpenDirectoryDialog(a.Context(), rpkg.OpenDialogOptions{})
	if err != nil {
		return models.ProtoDir{}, nil
	}

	protoFiles, err := filesystem.SearchProtoFiles([]string{path})
	if err != nil {
		if errors.Is(err, models.ErrNoProtoFilesFound) {
			a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
				Type:    "error",
				Title:   "Can't create a workspace",
				Message: "No proto files found",
			})
		}
		return models.ProtoDir{}, fmt.Errorf("api: failed to search proto files: %w", err)
	}

	return models.ProtoDir{
		Dir:   path,
		Files: protoFiles.RelativeProtoPaths,
	}, nil
}

func (a *Api) CreateWorkspace(name string, dirs []string) (models.Workspace, error) {
	registry, err := a.protoRegistryFromPath(dirs)
	if err != nil {
		var e reporter.ErrorWithPos
		if errors.As(err, &e) {
			var pathE *fs.PathError
			if errors.As(e.Unwrap(), &pathE) {
				pos := e.GetPosition()
				a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
					Type:    "error",
					Title:   fmt.Sprintf("Can't resolve import proto file %s", pathE.Path),
					Message: fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Col),
				})
			}

			if strings.Contains(e.Error(), "already defined at") {
				pos := e.GetPosition()
				a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
					Type:    "error",
					Title:   "Duplicated type definition found",
					Message: fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Col),
				})
			}
		}
		if errors.Is(err, models.ErrNoProtoFilesFound) {
			a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
				Type:    "error",
				Title:   "Can't create a workspace",
				Message: "No proto files found",
			})
		}
		return models.Workspace{}, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	spc, err := registry.Schema()
	if err != nil {
		return models.Workspace{}, fmt.Errorf("api: failed to create spec from registry: %w", err)
	}
	if len(spc.Services) == 0 {
		a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
			Type:    "error",
			Title:   "Can't create a workspace",
			Message: "No services found",
		})
		return models.Workspace{}, fmt.Errorf("no services found")
	}

	ws := models.Workspace{
		ID:          uuid.NewString(),
		Name:        name,
		Spec:        spc,
		BasePath:    dirs,
		TargetUrl:   "localhost:9000",
		LastUsage:   time.Now().UTC().Round(time.Nanosecond),
		ScriptFiles: make([]models.File, 0),
	}
	if err := a.store.SaveWorkspace(ws); err != nil {
		return ws, fmt.Errorf("api: failed to save workspace: %w", err)
	}

	a.registryStore.Add(ws.ID, registry)

	return ws, nil
}

func (a *Api) CreateWorkspaceV2(name string, dirs []string, workspaceKind models.WorkspaceKind) (models.Workspace, error) {
	var registry definitions.Registry
	var err error

	switch workspaceKind {
	case models.WorkspaceKindProto:
		registry, err = a.protoRegistryFromPath(dirs)
	case models.WorkspaceKindOpenapi:
		registry, err = a.openapiRegistryFromPath(dirs)
	default:
		return models.Workspace{}, fmt.Errorf("given unknown workspace kind: %s", workspaceKind)
	}

	if err != nil {
		var e reporter.ErrorWithPos
		if errors.As(err, &e) {
			var pathE *fs.PathError
			if errors.As(e.Unwrap(), &pathE) {
				pos := e.GetPosition()
				a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
					Type:    "error",
					Title:   fmt.Sprintf("Can't resolve import proto file %s", pathE.Path),
					Message: fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Col),
				})
			}

			if strings.Contains(e.Error(), "already defined at") {
				pos := e.GetPosition()
				a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
					Type:    "error",
					Title:   "Duplicated type definition found",
					Message: fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Col),
				})
			}
		}
		if errors.Is(err, models.ErrNoProtoFilesFound) {
			a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
				Type:    "error",
				Title:   "Can't create a workspace",
				Message: "No proto files found",
			})
		}
		return models.Workspace{}, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	spc, err := registry.Schema()
	if err != nil {
		return models.Workspace{}, fmt.Errorf("api: failed to create spec from registry: %w", err)
	}
	if len(spc.Services) == 0 {
		a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
			Type:    "error",
			Title:   "Can't create a workspace",
			Message: "No services found",
		})
		return models.Workspace{}, fmt.Errorf("no services found")
	}

	ws := models.Workspace{
		ID:        uuid.NewString(),
		Name:      name,
		Spec:      spc,
		BasePath:  dirs,
		TargetUrl: "localhost:9000",
		LastUsage: time.Now().UTC().Round(time.Nanosecond),
		ScriptFiles: make([]models.File, 0),
	}
	if err := a.store.SaveWorkspace(ws); err != nil {
		return ws, fmt.Errorf("api: failed to save workspace: %w", err)
	}

	a.registryStore.Add(ws.ID, registry)

	return ws, nil
}

func (s *Api) DeleteWorkspace(id string) error {
	return s.store.DeleteWorkspace(id)
}

func (s *Api) UpdateWorkspace(ws models.Workspace) error {
	return s.store.SaveWorkspace(ws)
}

func (s *Api) RenameWorkspace(id, name string) error {
	ws, err := s.store.GetWorkspace(id)
	if err != nil {
		return fmt.Errorf("failed to get a ws %s: %w", id, err)
	}

	ws.Name = name
	return s.store.SaveWorkspace(ws)
}

func (a *Api) WorkspaceList(id string) (models.WorkspaceList, error) {
	var res models.WorkspaceList
	list, err := a.store.GetWorkspaces()
	if err != nil {
		return res, err
	}
	if len(list) == 0 {
		return res, nil
	}
	if id == "" {
		sort.Slice(list, func(i, j int) bool {
			return list[i].LastUsage.After(list[j].LastUsage)
		})
		id = list[0].ID
	}

	var main models.Workspace
	for i := range list {
		if list[i].ID == id {
			list[i].LastUsage = time.Now().UTC().Round(time.Nanosecond)
			main = list[i]
			break
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].LastUsage.After(list[j].LastUsage)
	})

	main, err = a.enrichWorkspace(main)
	if err != nil {
		var ePos reporter.ErrorWithPos
		if errors.As(err, &ePos) {
			var pathE *fs.PathError
			if errors.As(ePos.Unwrap(), &pathE) {
				pos := ePos.GetPosition()
				a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
					Type:    "error",
					Title:   fmt.Sprintf("Can't resolve import proto file %s", pathE.Path),
					Message: fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Col),
				})
			}
		}

		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			click, _ := a.runtime.MessageDialog(a.ctx, rpkg.MessageDialogOptions{
				Type:         rpkg.QuestionDialog,
				Title:        fmt.Sprintf("Workspace '%s' can't start", main.Name),
				Message:      fmt.Sprintf("%s: no such file or directory.\nDelete the workspace?", main.BasePath),
				Buttons:      []string{"Yes", "No"},
				CancelButton: "No",
			})
			if click == "Yes" {
				a.DeleteWorkspace(main.ID)
			}
		}
		return res, err

	}

	shortList := make([]models.WorkspaceShort, len(list))
	for i := range list {
		shortList[i] = models.WorkspaceShort{
			ID:   list[i].ID,
			Name: list[i].Name,
		}
	}

	res.List = shortList
	res.Main = main
	return res, nil
}

func (s *Api) enrichWorkspace(ws models.Workspace) (models.Workspace, error) {
	registry, err := s.protoRegistryFromPath(ws.BasePath)
	if err != nil {
		return ws, err
	}
	s.registryStore.Add(ws.ID, registry)

	spec, err := registry.Schema()
	if err != nil {
		return ws, err
	}

	newWs := ws
	newWs.Spec = spec

	if !reflect.DeepEqual(ws, newWs) {
		if err := s.store.SaveWorkspace(newWs); err != nil {
			return newWs, err
		}
	}

	return newWs, nil
}

// ENVIRONMENT API

func (s *Api) GetGlobalVars() (string, error) {
	return s.store.GlobalVars()
}

func (s *Api) SaveGlovalVars(vars string) error {
	if err := s.store.SaveGlobalVars(vars); err != nil {
		return err
	}
	ip := interpreter.NewInterpreter("")
	if _, err := ip.Raw(vars); err != nil {
		return err
	}

	return nil
}

// GRPC API

func (a *Api) SendGrpc(request models.Request) (models.Response, error) {
	if strings.TrimSpace(request.Body) == "" {
		return models.Response{}, nil
	}

	reg, err := a.registryStore.Get(request.WorkspaceID)
	if err != nil {
		return models.Response{}, err
	}

	ws, err := a.store.GetWorkspace(request.WorkspaceID)
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

	vars, err := a.GetGlobalVars()
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to get global vars: %w", err)
	}
	ip := interpreter.NewInterpreter(vars)
	specInputMessage, err := ws.Spec.FindInputMessage(models.MethodName(request.Method).ServiceAndShort())
	if err != nil {
		return models.Response{}, err
	}

	inputType, err := reg.GetInputType(request.Method)
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to find input type: %w", err)
	}
	req, err := ip.CreateMessageFromScript(request.Body, inputType, ws.Spec, specInputMessage)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create request: %w", err)
	}

	meta, err := ip.CreateMetadata(request.Meta)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to create metadata: %w", err)
	}

	outputType, err := reg.GetOutputType(request.Method)
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to find output type: %w", err)
	}
	resp := dynamic.NewMessage(outputType)

	ctx = metadata.NewOutgoingContext(ctx, meta)
	responseMeta := metadata.MD{}
	path, err := reg.MethodPath(request.Method)
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to find method path: %w", err)
	}
	apiErr, err := c.Invoke(ctx, path, req, resp, &responseMeta)
	if err != nil {
		return models.Response{}, fmt.Errorf("api: failed to invoke method: %w", err)
	}

	specOutputMessage, err := ws.Spec.FindOutputMessage(models.MethodName(request.Method).ServiceAndShort())
	if err != nil {
		return models.Response{}, err
	}
	var body string
	if apiErr != "" {
		body = apiErr
	} else {
		body, err = render.New(reg.Links()).MessageAsJsString(specOutputMessage, resp)
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

// SCRIPTING API

func (a *Api) RunScript(request models.ScriptCall) (string, error) {
	if strings.TrimSpace(request.Body) == "" {
		return "", nil
	}

	reg, err := a.registryStore.Get(request.WorkspaceID)
	if err != nil {
		return "", err
	}

	ws, err := a.store.GetWorkspace(request.WorkspaceID)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	c, err := a.newClient(ctx, request.Addr)
	if err != nil {
		return "", fmt.Errorf("api: failed to create client: %w", err)
	}
	defer func() { _ = c.Close() }()

	vars, err := a.GetGlobalVars()
	if err != nil {
		return "", fmt.Errorf("failed to get global vars: %w", err)
	}
	ip := interpreter.NewInterpreter(vars)

	resp, err := ip.RunScript(ctx, request.Body, request.Meta, ws.Spec, reg, c, render.New(reg.Links()))
	if err != nil {
		return "", fmt.Errorf("api: failed to create request: %w", err)
	}
	return resp, nil
}

func (s *Api) protoRegistryFromPath(dirs []string) (*pcompiler.Registry, error) {
	protoFiles, err := filesystem.SearchProtoFiles(dirs)
	if err != nil {
		return nil, fmt.Errorf("api: failed to search proto files: %w", err)
	}

	registry, err := s.protoCompiler.Compile(protoFiles.AbsoluteDirsPath, protoFiles.RelativeProtoPaths, protoFiles.BufDirs)
	if err != nil {
		return nil, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	return registry, nil
}

func (s *Api) openapiRegistryFromPath(dirs []string) (*ocompiler.Registry, error) {
	protoFiles, err := filesystem.SearchOpenapiFiles(dirs)
	if err != nil {
		return nil, fmt.Errorf("api: failed to search openapi files: %w", err)
	}

	registry, err := s.openapiCompiler.Compile(protoFiles)
	if err != nil {
		return nil, fmt.Errorf("api: failed to compile proto files: %w", err)
	}

	return registry, nil
}

// SCRIPTING FILES API

func (s *Api) CreateScriptFile(workspaceID, name, content string) (models.File, error) {
	file := models.File{
		Id:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now().UTC().Round(time.Nanosecond),
		Content:   content,
	}

	ws, err := s.store.GetWorkspace(workspaceID)
	if err != nil {
		return file, err
	}

	ws.ScriptFiles = append([]models.File{file}, ws.ScriptFiles...)
	err = s.store.SaveWorkspace(ws)
	return file, err
}

func (s *Api) RemoveScriptFile(workspaceID, fileID string) ([]models.File, error) {
	ws, err := s.store.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.File, 0, len(ws.ScriptFiles))
	for _, file := range ws.ScriptFiles {
		if file.Id == fileID {
			continue
		}

		filtered = append(filtered, file)
	}

	ws.ScriptFiles = filtered
	err = s.store.SaveWorkspace(ws)
	return ws.ScriptFiles, err
}

func (s *Api) RenameScriptFile(workspaceID, fileID, name string) error {
	ws, err := s.store.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}

	for i, file := range ws.ScriptFiles {
		if file.Id == fileID {
			ws.ScriptFiles[i].Name = name
			break
		}
	}

	return s.store.SaveWorkspace(ws)
}

func (s *Api) UpdateScriptFileContent(workspaceID, fileID, content string) error {
	ws, err := s.store.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}

	for i, file := range ws.ScriptFiles {
		if file.Id == fileID {
			ws.ScriptFiles[i].Content = content
			break
		}
	}

	return s.store.SaveWorkspace(ws)
}
