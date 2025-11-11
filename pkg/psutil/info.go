package psutil

import (
	"runtime"

	"github.com/gitmanntoo/tinytop/pkg/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

type GoInfo struct {
	GOOS    string `json:"go_os"`
	GOARCH  string `json:"go_arch"`
	NumCPU  int    `json:"num_cpu"`
	Version string `json:"version"`
}

type PartitionUsage struct {
	Partition disk.PartitionStat `json:"partition"`
	Usage     *disk.UsageStat    `json:"usage"`
	UsageErr  error              `json:"usage_err,omitempty"`
}

type SysInfo struct {
	// Go runtime
	GoInfo GoInfo `json:"go_info"`
	// gopsutil
	// cpu
	CPUCores int            `json:"cpu_cores"`
	CPUInfo  []cpu.InfoStat `json:"cpu_info"`
	// disk
	IOCounters map[string]disk.IOCountersStat `json:"io_counters"`
	Partitions []PartitionUsage               `json:"partitions"`
	// host
	HostInfo *host.InfoStat `json:"host_info"`
	// mem
	MemInfo     *mem.VirtualMemoryStat `json:"mem_info"`
	SwapDevices []*mem.SwapDevice       `json:"swap_devices,omitempty"`
	SwapInfo    *mem.SwapMemoryStat    `json:"swap_info"`
	// sensors
	SensorTemperatures []sensors.TemperatureStat `json:"sensor_temperatures,omitempty"`
}

// Info retrieves static system information.
func Info() (*SysInfo, error) {
	info := &SysInfo{
		// Go runtime info
		GoInfo: GoInfo{
			GOOS:    runtime.GOOS,
			GOARCH:  runtime.GOARCH,
			NumCPU:  runtime.NumCPU(),
			Version: runtime.Version(),
		},
	}

	// Get CPU core count and info
	if counts, err := cpu.Counts(true); err != nil {
		return nil, err
	} else {
		info.CPUCores = counts
	}

	if stats, err := cpu.Info(); err != nil {
		return nil, err
	} else {
		info.CPUInfo = stats
	}

	// Get disk partitions
	parts, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	// Get usage for each partition
	info.Partitions = []PartitionUsage{}
	for _, part := range parts {
		// Exclude partitions with "nobrowse" option (macOS)
		var isNoBrowse bool
		for _, opt := range part.Opts {
			if opt == "nobrowse" {
				isNoBrowse = true
				break
			}
		}
		if isNoBrowse {
			continue
		}

		usage, err := disk.Usage(part.Mountpoint)
		pu := PartitionUsage{
			Partition: part,
			Usage:     usage,
			UsageErr:  err,
		}
		info.Partitions = append(info.Partitions, pu)
	}

	// Get disk IO counters
	if ioCounters, err := disk.IOCounters(); err != nil {
		return nil, err
	} else {
		info.IOCounters = ioCounters
	}

	// Get system host information
	if hostInfo, err := host.Info(); err != nil {
		return nil, err
	} else {
		info.HostInfo = hostInfo
	}

	// Get memory information
	if memInfo, err := mem.VirtualMemory(); err != nil {
		return nil, err
	} else {
		info.MemInfo = memInfo
	}

	if swapDevices, err := mem.SwapDevices(); err != nil {
		utils.Log.Warn().Err(err).Msg("Failed to get swap devices")
	} else {
		info.SwapDevices = swapDevices
	}

	if swapInfo, err := mem.SwapMemory(); err != nil {
		return nil, err
	} else {
		info.SwapInfo = swapInfo
	}

	// Get sensor temperatures
	if temps, err := sensors.SensorsTemperatures(); err != nil {
		utils.Log.Warn().Err(err).Msg("Failed to get sensor temperatures")
	} else {
		info.SensorTemperatures = temps
	}

	return info, nil
}
