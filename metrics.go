package sdk

import (
	"log"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func getSystemMetrics() (cpuPercent float64, memTotal, memAvail int64) {
	if percents, err := cpu.Percent(0, false); err == nil && len(percents) > 0 {
		cpuPercent = percents[0]
	} else if err != nil {
		log.Printf("[Anox SDK] Failed to read CPU metrics: %v", err)
	}
	if vm, err := mem.VirtualMemory(); err == nil {
		memTotal = int64(vm.Total / 1024 / 1024)
		memAvail = int64(vm.Available / 1024 / 1024)
	} else if err != nil {
		log.Printf("[Anox SDK] Failed to read memory metrics: %v", err)
	}
	return cpuPercent, memTotal, memAvail
}
