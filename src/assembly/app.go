package assembly

import (
	"context"
	"kalisto/src/api"
	"kalisto/src/db"
	"kalisto/src/environment"
	"kalisto/src/proto/client"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"
	"log"

	"github.com/adrg/xdg"
)

// App struct
type App struct {
	ctx context.Context

	Api *api.Api
}

// NewApp creates a new App application struct
func NewApp() *App {
	store, err := db.New(xdg.DataHome)
	if err != nil {
		log.Fatal("failed to init db: ", err)
	}

	protoRegistry := compiler.NewProtoRegistry()
	ws, err := workspace.New(store)
	if err != nil {
		log.Fatal("failed to init workspace: ", err)
	}
	glovalVars, err := environment.NewGlovalVars(store)
	if err != nil {
		log.Fatal("failed to init environments: ", err)
	}
	protoCompiler := compiler.NewFileCompiler()
	specFactory := spec.NewFactory()

	newClient := func(ctx context.Context, addr string) (*client.Client, error) {
		return client.NewClient(ctx, client.Config{
			Addr: addr,
		})
	}

	a := api.New(protoCompiler, specFactory, ws, glovalVars, newClient, protoRegistry)

	return &App{Api: a}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Start(ctx context.Context) {
	a.ctx = ctx

	api.SetContext(a.Api, ctx)
}
