package main

import (
	"context"
)

type App struct {
	ctx context.Context
}

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
}
