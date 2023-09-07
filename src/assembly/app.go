package assembly

import (
	"context"
	"fmt"
	"kalisto/src/api"
	"kalisto/src/db"
	"kalisto/src/pkg/runtime"
	"kalisto/src/proto/client"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/spec"

	"github.com/adrg/xdg"
)

// App struct
type App struct {
	ctx context.Context

	Api *api.Api

	db *db.DB
}

// NewApp creates a new App application struct
func NewApp() (*App, error) {
	store, err := db.New(xdg.DataHome)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	protoRegistry := compiler.NewProtoRegistry()
	protoCompiler := compiler.NewFileCompiler()
	specFactory := spec.NewFactory()

	newClient := func(ctx context.Context, addr string) (*client.Client, error) {
		return client.NewClient(ctx, client.Config{
			Addr: addr,
		})
	}

	a := api.New(protoCompiler, specFactory, store, newClient, protoRegistry, runtime.New())

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
