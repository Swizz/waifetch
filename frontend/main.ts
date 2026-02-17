import { effect, sprae, store } from "sprae"
import { GetSystemInfo } from "./wailsjs/go/main/App.js"

interface SystemData {
	user: string
	hostname: string
	os: string
	platform: string
	kernel: string
	cpu: string
	uptime: number
	disk: { used: number; total: number }
	mem: { used: number; total: number }
	dark: boolean
}

interface DataFormatter {
	_fmtUptime(uptime: number): string
	_fmtBytes(bytes: number): string
	_fmtPercent(value: number): string
	_platformImg(os: string, platform: string, path?: string): string
}

const emptyState: SystemData = {
	user: "",
	hostname: "",
	os: "",
	platform: "",
	kernel: "",
	cpu: "",
	uptime: 0,
	disk: { used: 0, total: 0 },
	mem: { used: 0, total: 0 },
	dark: false,
}

const formatter: DataFormatter = {
	_fmtUptime(uptime) {
		// @ts-ignore: Intl.DurationFormat is a newer proposal
		return new Intl.DurationFormat("en", { style: "narrow" }).format({
			days: Math.floor(uptime / 60 / 60 / 24),
			hours: Math.floor(uptime / 60 / 60) % 24,
			minutes: Math.floor(uptime / 60) % 60,
			seconds: uptime % 60,
		})
	},

	_fmtBytes(bytes) {
		const units = ["byte", "kilobyte", "megabyte", "gigabyte", "terabyte"]
		const n = bytes && Math.floor(Math.log(bytes) / Math.log(1024))
		return new Intl.NumberFormat("en", {
			style: "unit",
			unit: units[n],
			unitDisplay: "narrow",
			maximumFractionDigits: 2,
		}).format(bytes / Math.pow(1024, n))
	},

	_fmtPercent(value) {
		return new Intl.NumberFormat("en", {
			style: "percent",
			unitDisplay: "narrow",
			maximumFractionDigits: 0,
		}).format(value || 0)
	},

	_platformImg(
		os,
		platform,
		path = "./loganmarchione/homelab-svg-assets/{name}.svg",
	) {
		if (os == "linux") {
			if (platform == "arch") return path.replace("{name}", platform + os)
			if (platform == "ubuntu") {
				return path.replace("{name}", "canonical" + platform)
			}
		}

		if (os == "windows") {
			const match = new RegExp(/windows +([0-9]+)/i).exec(platform)
			if (match) return path.replace("{name}", os + match[1])
		}

		return path.replace("{name}", platform)
	},
}

const state: SystemData & DataFormatter = store({ ...emptyState, ...formatter })

effect(() => {
	GetSystemInfo().then((systemData: SystemData) =>
		Object.assign(state, systemData)
	)
})

export default sprae(document.body, state)
