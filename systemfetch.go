package main

import (
	"fmt"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	wails "github.com/wailsapp/wails/v3/pkg/application"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/thiagokokada/dark-mode-go"
)

type SystemFetch struct{}

type UsageInfo struct {
	Used  uint64 `json:"used,omitempty"`
	Total uint64 `json:"total,omitempty"`
}

type SystemInfo struct {
	User     string     `json:"user,omitempty"`
	Hostname string     `json:"hostname,omitempty"`
	OS       string     `json:"os,omitempty"`
	Platform string     `json:"platform,omitempty"`
	Kernel   string     `json:"kernel,omitempty"`
	Cpu      string     `json:"cpu,omitempty"`
	Uptime   uint64     `json:"uptime,omitempty"`
	Disk     *UsageInfo `json:"disk,omitempty"`
	Mem      *UsageInfo `json:"mem,omitempty"`
	Dark     bool       `json:"dark"`
}

type Event string

const (
	SystemInfoUpdate Event = "systeminfo:update"
)

func (_ *SystemFetch) GetSystemInfo() SystemInfo {
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
		system.Disk = &UsageInfo{}
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
		system.Mem = &UsageInfo{
			Total: memory.Total,
			Used:  memory.Used,
		}
	}

	var darkm bool
	darkm, err = dark.IsDarkMode()

	if err == nil {
		system.Dark = darkm
	}

	return system
}

func (_ *SystemFetch) MonitorSystemInfo() Event {
	var ticker *time.Ticker = time.NewTicker(3 * time.Second)
	var system SystemInfo = SystemInfo{}

	var app *wails.App = application.Get()

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			var err error

			var memory *mem.VirtualMemoryStat
			memory, err = mem.VirtualMemory()
			if err == nil {
				system.Mem = &UsageInfo{
					Total: memory.Total,
					Used:  memory.Used,
				}
			}

			var darkm bool
			darkm, err = dark.IsDarkMode()

			if err == nil {
				system.Dark = darkm
			}

			app.Event.Emit(string(SystemInfoUpdate), system)
		}
	}()

	return SystemInfoUpdate
}

func init() {
	application.RegisterEvent[SystemInfo](string(SystemInfoUpdate))
}
