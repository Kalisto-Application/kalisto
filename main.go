package main

import (
	"embed"
	"log"
	"time"

	"kalisto/src/assembly"
	"kalisto/src/config"
	"kalisto/src/models"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/getsentry/sentry-go"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	sentryDsn := config.C.SentryDsn
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: sentryDsn,
	}); err != nil {
		log.Fatalln("failed to initialize sentry client:", err)
	}
	defer sentry.Flush(time.Second * 3)
	defer sentry.Recover()

	// Create an instance of the app structure
	app, err := assembly.NewApp()
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	// Create application with options
	if err := wails.Run(&options.App{
		Title:  "kalisto",
		Width:  1440,
		Height: 925,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Start,
		Bind: []interface{}{
			app.Api,
		},
		ErrorFormatter: models.NewErrorFormatter(func(err error) {
			sentry.CaptureException(err)
		}),
	}); err != nil {
		sentry.CaptureException(err)
	}
}
