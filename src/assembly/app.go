package assembly

import (
	"context"
	"fmt"
	"kalisto/src/api"
	"kalisto/src/db"
	"kalisto/src/definitions"
	ocompiler "kalisto/src/definitions/openapi/compiler"
	"kalisto/src/definitions/proto/client"
	pcompiler "kalisto/src/definitions/proto/compiler"
	"kalisto/src/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	Api *api.Api

	db *db.DB
}

// NewApp creates a new App application struct
func NewApp(homeDir string) (*App, error) {
	store, err := db.New(homeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	protoCompiler := pcompiler.NewFileCompiler()
	openapiCompiler := ocompiler.NewCompiler()

	newClient := func(ctx context.Context, addr string) (*client.Client, error) {
		return client.NewClient(ctx, client.Config{
			Addr: addr,
		})
	}

	a := api.New(protoCompiler, openapiCompiler, store, newClient, definitions.NewRegistryStore(), runtime.New())

	return &App{Api: a, db: store}, nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Start(ctx context.Context) {
	a.ctx = ctx

	api.SetContext(a.Api, ctx)
}

func (a *App) OnShutdown(ctx context.Context) {
	a.db.Close()
}
