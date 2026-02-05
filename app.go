package main

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/thiagokokada/dark-mode-go"
)

type UsageInfo struct {
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
}

type SystemInfo struct {
	User     string    `json:"user"`
	Hostname string    `json:"hostname"`
	OS       string    `json:"os"`
	Platform string    `json:"platform"`
	Kernel   string    `json:"kernel"`
	Cpu      string    `json:"cpu"`
	Uptime   uint64    `json:"uptime"`
	Disk     UsageInfo `json:"disk"`
	Mem      UsageInfo `json:"mem"`
	Dark     bool      `json:"dark"`
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

	var cpus []cpu.InfoStat
	cpus, err = cpu.Info()

	if err == nil && len(cpus) > 0 {
		system.Cpu = cpus[0].ModelName
	}

	var partitions []disk.PartitionStat
	partitions, err = disk.Partitions(true)

	if err == nil {
		var partition disk.PartitionStat
		for _, partition = range partitions {
			var disku *disk.UsageStat
			disku, err = disk.Usage(partition.Mountpoint)

			if err == nil {
				system.Disk.Total += disku.Total
				system.Disk.Used += disku.Used
			}
		}
	}

	var memory *mem.VirtualMemoryStat
	memory, err = mem.VirtualMemory()

	if err == nil {
		system.Mem.Total = memory.Total
		system.Mem.Used = memory.Used
	}

	var darkm bool
	darkm, err = dark.IsDarkMode()

	if err == nil {
		system.Dark = darkm
	}

	return system
}
