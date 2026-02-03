package main

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/host"
)

type SystemInfo struct {
	User     string `json:"user"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Platform string `json:"platform"`
	Kernel   string `json:"kernel"`
	Uptime   uint64 `json:"uptime"`
}

type App struct {
	ctx context.Context
}

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
}

func (app *App) GetSystemInfo() SystemInfo {
	var err error
	var system SystemInfo = SystemInfo{}

	var info *host.InfoStat
	info, err = host.Info()

	if err == nil {
		system.Hostname = info.Hostname
		system.OS = info.OS
		system.Platform = info.Platform
		system.Kernel = fmt.Sprintf("%s %s", info.KernelArch, info.KernelVersion)
		system.Uptime = info.Uptime
	}

	var users []host.UserStat
	users, err = host.Users()

	if err == nil && len(users) > 0 {
		system.User = users[0].User
	}

	return system
}
