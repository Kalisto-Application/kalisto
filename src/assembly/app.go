package assembly

import (
	"context"
	"kalisto/src/api"
	"kalisto/src/db"
	"kalisto/src/environment"
	"kalisto/src/proto/compiler"
	"kalisto/src/proto/spec"
	"kalisto/src/workspace"
	"log"
)

// App struct
type App struct {
	ctx context.Context

	Api *api.Api
}

// NewApp creates a new App application struct
func NewApp() *App {
	store, err := db.New()
	if err != nil {
		log.Fatal("failed to init db: ", err)
	}

	ws, err := workspace.New(store)
	if err != nil {
		log.Fatal("failed to init workspace: ", err)
	}
	env, err := environment.NewEnvironment(store)
	if err != nil {
		log.Fatal("failed to init environments: ", err)
	}
	protoCompiler := compiler.NewFileCompiler()
	specFactory := spec.NewFactory()

	a := api.New(protoCompiler, specFactory, ws, env)

	return &App{Api: a}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Start(ctx context.Context) {
	a.ctx = ctx
}
