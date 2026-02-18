package main

import (
	"embed"

	wails "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var app *wails.App

	app = wails.New(wails.Options{
		Name: "waifetch",
		Services: []wails.Service{
			wails.NewService(&SystemFetch{}),
		},
		Assets: wails.AssetOptions{
			Handler: wails.AssetFileServerFS(assets),
		},
	})

	app.Window.NewWithOptions(wails.WebviewWindowOptions{
		Title:         "waifetch",
		Width:         1287,
		Height:        728,
		DisableResize: true,
	})

	var err error = app.Run()

	if err != nil {
		println("Error:", err.Error())
	}
}
