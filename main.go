package main

import (
	"embed"
	"fmt"
	"runtime"
	"time"

	"kalisto/src/assembly"
	"kalisto/src/config"
	"kalisto/src/models"
	"kalisto/src/pkg/log"
	"kalisto/src/pkg/update"
	stdlog "log"

	"github.com/adrg/xdg"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/getsentry/sentry-go"
)

var version string
var ghApiToken string
var platform string

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	l := log.New()
	sentryDsn := config.C.SentryDsn
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:     sentryDsn,
		Release: version,
	}); err != nil {
		stdlog.Fatalln("failed to initialize sentry client:", err)
	}
	defer sentry.Flush(time.Second * 3)
	defer sentry.Recover()
	l.Debug("sentry initialized")

	AppMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu())
	}
	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText(fmt.Sprintf("Version %s", version), nil, func(_ *menu.CallbackData) {})
	l.Debug("menu built")

	updater := update.NewUpdater(version, platform, ghApiToken)
	updated, err := updater.Run()
	if err != nil {
		l.Error(err.Error())
		sentry.CaptureException(err)
	}
	if updated {
		l.Info("app restart after update")
		if err := updater.Restart(); err != nil {
			l.Error(err.Error())
			sentry.CaptureException(err)
		}
	}
	l.Debug("updater ran")

	// Create an instance of the app structure
	app, err := assembly.NewApp(xdg.DataHome)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	// Create application with options
	if err := wails.Run(&options.App{
		OnShutdown: app.OnShutdown,
		Title:      "kalisto",
		Width:      1440,
		Height:     925,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Start,
		Bind: []interface{}{
			app.Api,
		},
		ErrorFormatter: models.NewErrorFormatter(app.Api.Context, func(err error) {
			sentry.CaptureException(err)
		}, app.Api.Runtime()),
		Menu:   AppMenu,
		Logger: l,
	}); err != nil {
		sentry.CaptureException(err)
	}
}
