package assembly

import (
	"context"
	"kalisto/src/api"
)

// App struct
type App struct {
	ctx context.Context

	Api *api.Api
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := api.New()

	return &App{Api: a}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Start(ctx context.Context) {
	a.ctx = ctx
}
