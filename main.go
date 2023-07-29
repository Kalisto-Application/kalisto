package main

import (
	"embed"
	"log"
	"time"

	"kalisto/src/assembly"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/getsentry/sentry-go"
)

//go:embed all:frontend/dist
var assets embed.FS
var sentryDsn string

func main() {
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
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Start,
		Bind: []interface{}{
			app.Api,
		},
		ErrorFormatter: func(err error) any {
			sentry.CaptureException(err)
			return err
		},
	}); err != nil {
		sentry.CaptureException(err)
	}
}
